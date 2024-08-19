# Cloud Compliance Framework - Template Plugin
# template-cf-plugin


## Development

First, change main.go to suit your needs.

Look for lines that have 'GITHUB TEMPLATE INSTRUCTIONS' in them, and follow them.

Then run

```sh
go mod init [YOUR PROVIDER NAME]-cf-plugin     # eg yourcloudprovider-cf-plugin
```

If you push to GitHub with an appropriate `GITHUB_TOKEN` in your secrets,
then the image should be built and made publicly-available to Compliance Framework.
