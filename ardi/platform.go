package ardi

import (
	"context"
	"io"
	"strings"

	rpc "github.com/arduino/arduino-cli/rpc/commands"
)

func updateLibraryIndex(client rpc.ArduinoCoreClient, instance *rpc.Instance) {
	logger.Debug("Updating library index")
	libIdxUpdateStream, err := client.UpdateLibrariesIndex(context.Background(),
		&rpc.UpdateLibrariesIndexReq{Instance: instance})

	if err != nil {
		logger.WithError(err).Fatal("Error updating libraries index")
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		resp, err := libIdxUpdateStream.Recv()
		if err == io.EOF {
			logger.Debug("Library index update done")
			break
		}

		if err != nil {
			logger.WithError(err).Fatal("Error updating libraries index")
		}

		if resp.GetDownloadProgress() != nil {
			logger.Debugf("DOWNLOAD: %s", resp.GetDownloadProgress())
		}
	}
}

func updatePlatformIndex(client rpc.ArduinoCoreClient, instance *rpc.Instance) {
	logger.Debug("Updating index...")
	uiRespStream, err := client.UpdateIndex(context.Background(), &rpc.UpdateIndexReq{
		Instance: instance,
	})
	if err != nil {
		logger.Fatalf("Error updating index: %s", err)
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		uiResp, err := uiRespStream.Recv()

		// the server is done
		if err == io.EOF {
			logger.Debug("Update index done")
			break
		}

		// there was an error
		if err != nil {
			logger.Fatalf("Update error: %s", err)
		}

		// operations in progress
		if uiResp.GetDownloadProgress() != nil {
			logger.Debugf("DOWNLOAD: %s", uiResp.GetDownloadProgress())
		}
	}
}

func platformUpgrade(client rpc.ArduinoCoreClient, instance *rpc.Instance, platPackage, arch string, done chan platformUpgradeMessage) {
	logger.Debugf("Upgrading platform: %s:%s\n", platPackage, arch)

	upgradeRespStream, err := client.PlatformUpgrade(context.Background(),
		&rpc.PlatformUpgradeReq{
			Instance:        instance,
			PlatformPackage: platPackage,
			Architecture:    arch,
		})

	if err != nil {
		logger.WithError(err).Warn("Error upgrading platform")
	}

	message := platformUpgradeMessage{
		platformPackage: platPackage,
		architecture:    arch,
		success:         false,
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		upgradeResp, err := upgradeRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			logger.Debug("Upgrade done")
			message.success = true
			done <- message
			break
		}

		// There was an error.
		if err != nil {
			if !strings.Contains(err.Error(), "platform already at latest version") {
				logger.WithError(err).Warn("Cannot upgrade platform")
			}
			done <- message
			break
		}

		// When a download is ongoing, log the progress
		if upgradeResp.GetProgress() != nil {
			logger.Debugf("DOWNLOAD: %s", upgradeResp.GetProgress())
		}

		// When an overall task is ongoing, log the progress
		if upgradeResp.GetTaskProgress() != nil {
			logger.Debugf("TASK: %s", upgradeResp.GetTaskProgress())
		}
	}
}

func upgradePlatforms(client rpc.ArduinoCoreClient, instance *rpc.Instance, platforms []*rpc.Platform) {
	done := make(chan platformUpgradeMessage)
	waitForAllJobs := make(chan bool)
	goRoutineSlot := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		goRoutineSlot <- struct{}{}
	}
	go func() {
		for i := 0; i < len(platforms); i++ {
			message := <-done
			if message.success {
				logger.Debugf("Successfully upgraded %s:%s", message.platformPackage, message.architecture)
			}
			// job has finished, release the go routine slot so another job can start
			goRoutineSlot <- struct{}{}
		}
		// signal all jobs complete
		waitForAllJobs <- true
	}()
	for _, plat := range platforms {
		// Wait for an available go routine slot before beginning the job
		<-goRoutineSlot
		id := plat.GetID()
		idParts := strings.Split(id, ":")
		platPackage := idParts[0]
		arch := idParts[len(idParts)-1]
		go platformUpgrade(client, instance, platPackage, arch, done)
	}

	<-waitForAllJobs
}

func platformInstall(client rpc.ArduinoCoreClient, instance *rpc.Instance, platPackage, arch, version string, done chan platformInstallMessage) {
	logger.Debugf("Installing platform: %s:%s\n", arch, version)

	installRespStream, err := client.PlatformInstall(
		context.Background(),
		&rpc.PlatformInstallReq{
			Instance:        instance,
			PlatformPackage: platPackage,
			Architecture:    arch,
			Version:         version,
		})

	if err != nil {
		logger.WithError(err).Warn("Failed to install platform")
	}

	message := platformInstallMessage{
		platformPackage: platPackage,
		architecture:    arch,
		version:         version,
		success:         false,
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		installResp, err := installRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			logger.Debug("Install done")
			message.success = true
			done <- message
			break
		}

		// There was an error.
		if err != nil {
			logger.WithError(err).Warn("Failed to install platform")
			done <- message
			break
		}

		// When a download is ongoing, log the progress
		if installResp.GetProgress() != nil {
			logger.Debugf("DOWNLOAD: %s", installResp.GetProgress())
		}

		// When an overall task is ongoing, log the progress
		if installResp.GetTaskProgress() != nil {
			logger.Debugf("TASK: %s", installResp.GetTaskProgress())
		}
	}
}

func loadAllPlatforms(client rpc.ArduinoCoreClient, instance *rpc.Instance) {
	searchResp, err := client.PlatformSearch(context.Background(), &rpc.PlatformSearchReq{
		Instance: instance,
	})

	if err != nil {
		logger.Fatalf("Search error: %s", err)
	}

	platforms := searchResp.GetSearchOutput()
	done := make(chan platformInstallMessage)
	waitForAllJobs := make(chan bool)
	goRoutineSlot := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		goRoutineSlot <- struct{}{}
	}
	go func() {
		for i := 0; i < len(platforms); i++ {
			message := <-done
			if message.success {
				logger.Debugf("Successfully installed %s:%s - %s", message.platformPackage, message.architecture, message.version)

			}
			// job has finished, release the go routine slot so another job can start
			goRoutineSlot <- struct{}{}
		}
		// signal all jobs complete
		waitForAllJobs <- true
	}()
	for _, plat := range platforms {
		// Wait for an available go routine slot before beginning the job
		<-goRoutineSlot
		id := plat.GetID()
		idParts := strings.Split(id, ":")
		platPackage := idParts[0]
		arch := idParts[len(idParts)-1]
		latest := plat.GetLatest()
		logger.Debugf("Search result: %s: %s - %s", platPackage, id, latest)
		go platformInstall(client, instance, platPackage, arch, latest, done)
	}
	<-waitForAllJobs
	upgradePlatforms(client, instance, platforms)
}

func platformList(client rpc.ArduinoCoreClient, instance *rpc.Instance) {
	listResp, err := client.PlatformList(context.Background(),
		&rpc.PlatformListReq{Instance: instance})

	if err != nil {
		logger.Fatalf("List error: %s", err)
	}

	logger.Debug("------INSTALLED PLATFORMS------")
	for _, plat := range listResp.GetInstalledPlatform() {
		logger.Debugf("Installed platform: %s - %s", plat.GetID(), plat.GetInstalled())
	}
	logger.Debug("-------------------------------")
}
