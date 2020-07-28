const core = require('@actions/core');
const { GitHub } = require('@actions/github');
const fs = require('fs');

async function run() {
  try {
    // Get authenticated GitHub client (Ocktokit): https://github.com/actions/toolkit/tree/master/packages/github#usage
    const github = new GitHub(process.env.GITHUB_TOKEN);

    const url = core.getInput('upload_url', { required: true });
    console.log({ upload_url: url });

    // Determine content-length for header to upload asset
    const contentLength = (filePath) => fs.statSync(filePath).size;
    const contentType = 'application/zip';
    const buildDir = 'v2/build';

    for (const file of fs.readdirSync(buildDir)) {
      const filePath = `${buildDir}/${file}`;
      const headers = {
        'content-type': contentType,
        'content-length': contentLength(filePath),
      };

      console.log('uploading', filePath);

      await github.repos.uploadReleaseAsset({
        url,
        headers,
        name: file,
        data: fs.createReadStream(filePath),
      });
    }
  } catch (e) {
    core.setFailed(e.message);
  }
}

run();
