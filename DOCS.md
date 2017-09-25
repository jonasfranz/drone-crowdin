Use the Crowdin plugin to update translation files.

You must provide in your configuration:

* `project_identifier` - Identifier of your Crowdin project (also available as secret  `CROWDIN_IDENTIFIER`)
* `project_key` - API Key of your Crowdin project (also available as secret  `CROWDIN_KEY`)
* `files` - Map of files to upload to Crowdin
  * key: the Crowdin file name
  * value: the real path the to file
* `ignore_branch` It will send the Drone branch to Crowdin if it is `false`. (Default: `false`)
* `download` Downloads translated files from Crowdin if it is `true`. (Default: `false`)
* `export_dir` Export directory of the translated strings
* `languages` Languages which should be downloaded/exported from Crowdin. (Default: `all`)


Information about API keys: https://support.crowdin.com/api/api-integration-setup/
## Example

The following is a sample configuration in your .drone.yml file:

```yaml
pipeline:
  crowdin:
    image: jonasfranz/crowdin
    project_identifier: example
    project_key: 1bc29b36f623ba82aaf6724fd3b16718
    files:
      example: options/example.ini
      example2: options/example2.ini
    ignore_branch: true
    download: true
    export_dir: langs/
    languages:
    - de
    - fr
```

## Commit changes

Please have a look at the [drone-git-push plugin](https://github.com/appleboy/drone-git-push) if you want to update the translations in your git repository too-

Example:
```yaml
pipeline:
  crowdin:
    image: jonasfranz/crowdin
    project_identifier: example
    project_key: 1bc29b36f623ba82aaf6724fd3b16718
    files:
      example: options/example.ini
      example2: options/example2.ini
    ignore_branch: true
    download: true
    export_dir: langs/
    languages:
    - de
    - fr
  git_push:
    image: appleboy/drone-git-push
    branch: master
    remote: git@your-remote.tdl/your-repo/repo
    force: false
    commit: true
    commit_message: "[skip ci] Updated translations"
```

**Important**: Please use `[skip ci]` inside your commit message to prevent recursive ci builds.
