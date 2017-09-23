Use the Crowdin plugin to update translation files.

You must provide in your configuration:

* `project_identifier` - Identifier of your Crowdin project
* `project_key` - API Key of your Crowdin project
* `files` - Map of files with the crowdin file name as key and the to-upload file-path as value.

Information about API keys: https://support.crowdin.com/api/api-integration-setup/
## Example

The following is a sample configuration in your .drone.yml file:

```yaml
translations:
  crowdin:
    project_identifier: example
    project_key: 1bc29b36f623ba82aaf6724fd3b16718
    files:
      example: options/example.ini
      example2: options/example2.ini
```