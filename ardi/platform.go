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

func platformUpgrade(client rpc.ArduinoCoreClient, instance *rpc.Instance, platPackage, arch string) {
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

	// Loop and consume the server stream until all the operations are done.
	for {
		upgradeResp, err := upgradeRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			logger.Debug("Upgrade done")
			break
		}

		// There was an error.
		if err != nil {
			if !strings.Contains(err.Error(), "platform already at latest version") {
				logger.WithError(err).Warn("Cannot upgrade platform")
			}
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

func platformInstall(client rpc.ArduinoCoreClient, instance *rpc.Instance, platPackage, arch, version string) {
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

	// Loop and consume the server stream until all the operations are done.
	for {
		installResp, err := installRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			logger.Debug("Install done")
			break
		}

		// There was an error.
		if err != nil {
			logger.WithError(err).Fatal("Failed to install platform")
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

	for _, plat := range platforms {
		id := plat.GetID()
		idParts := strings.Split(id, ":")
		platPackage := idParts[0]
		arch := idParts[len(idParts)-1]
		latest := plat.GetLatest()
		logger.Debugf("Search result: %s: %s - %s", platPackage, id, latest)
		platformInstall(client, instance, platPackage, arch, latest)
		platformUpgrade(client, instance, platPackage, arch)
	}
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
