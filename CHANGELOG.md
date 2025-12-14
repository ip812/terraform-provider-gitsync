## 1.3.0 (Dev 15, 2025)

ENHANCEMENTS:

* Add option to import all resources.

BUG FIXES:

* ci: Fix the wrong test commands.
* chore: Update dependencies to latest versions.

## 1.2.0 (Oct 27, 2025)

ENHANCEMENTS:

* resource/gitsync_values_yaml: Add content validation to ensure YAML format is correct.
* resource/gitsync_values_json: Add content validation to ensure JSON format is correct.

BUG FIXES:

* provider: Add exponential backoff retry mechanism when resource create/update fails(fixed race condition).
* ci: Remove unnecessary workflows.


## 1.1.0 (Oct 7, 2025)

ENHANCEMENTS:

* resource/gitsync_values_json: Add resource to manage a json values file in Git reposiotory. 
* resource/gitsync_values_file: Add resource to manage a values file(any type) in Git reposiotory. 

BUG FIXES:

* resource/gitsync_values_yaml: Add validation if the file's extension is .yaml or .yml. 


## 1.0.0 (Sep 28, 2025)

ENHANCEMENTS:

* resource/gitsync_values_yaml: Add resource to manage a yaml values file in Git reposiotory. 
