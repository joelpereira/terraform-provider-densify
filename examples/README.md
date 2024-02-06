# Examples

This directory contains examples that can be run/tested manually via the Terraform CLI.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or ar testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** cloud recommendation example for the provider
* **data-sources/\*** provider examples for pulling Densify cloud & container optimization recommendations as a Terraform data-source
