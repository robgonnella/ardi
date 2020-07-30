import * as core from '@actions/core';
import * as octo from '@actions/github';
import { createReadStream } from 'fs';
import { join, basename } from 'path';
import * as glob from 'fast-glob';
import * as mime from 'mime-types';

const repoFull = core.getInput('repo', { required: true });
const token = core.getInput('github_token', { required: true });
const releaseId = core.getInput('release_id', { required: true });
const filePattern = core.getInput('file_pattern', { required: true });
const github = octo.getOctokit(token);

const repoParts = repoFull.split('/');

const run = async function () {
  try {
    if (repoParts.length !== 2) {
      throw new Error(
        "Invalid repo value. Should be of form '<owner>/<repo_name>'",
      );
    }
    const owner = repoParts[0];
    const repo = repoParts[1];

    const id = parseInt(releaseId);
    if (isNaN(id)) {
      throw new Error(`invalid release_id: ${releaseId}`);
    }

    const files = await glob(filePattern.split(';'));

    for (const file of files) {
      const filePath = join(__dirname, file);
      const fileName = basename(filePath);
      const contentType = mime.lookup(filePath);
      if (!contentType) {
        throw new Error(`Unrecognized mime-type for file: ${filePath}`);
      }
      await github.repos.uploadReleaseAsset({
        owner,
        repo,
        release_id: id,
        name: fileName,
        data: '',
        file: createReadStream(filePath),
        headers: {
          'content-type': contentType,
        },
      });
    }
  } catch (error) {
    console.error(error);
    core.setFailed(error.message);
  }
};

run();
