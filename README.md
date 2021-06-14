# FileMango

An extensible framework for continuously creating and updating file metadata summaries using configurable external programs.

## Overview

Metadata is an incredibly valuable asset in the modern world. File metadata in particular can be leveraged by programs in order to improve a wide range of user experiences. In order to meet this demand for metadata, we developed and tested a framework for automatically mining and saving useful file metadata on a user’s system. The framework continuously watches for files of specific types within specific directories that are specified in a configuration file. The configuration file also contains a list of external programs that implement a standard file analyzer interface, and the file types that each program operates on. Once a file of the correct type is found, the framework runs one or many of the aforementioned external programs in parallel and specifies, through the standard interface, the file to analyze. Over time, the framework receives file metadata results from each external program. These results are then saved to the disk and associated directly with the file using extended attributes.

### Paper

A paper detailing our design and development process is located in paper/.

## How to Use

1. Download

Clone the repo or download a zip, no installation is necessary to use the executable. Compile dependencies are listed in [go.mod](https://github.com/Johnnydouble/FileMango/blob/master/go.mod).

2. Modules

FileMango relies on the specific programs you want to use, choose and install whatever programs you want to run on each file. 

3. Configure 

In order to run the analysis FileMango has to have a module interface program for each module you wish to use. Some examples are located [here](https://github.com/Johnnydouble/FileMangoExampleModule) and [here](https://github.com/holozene/exampleModule), fork and edit them or reverse engineer make your own (if you develop one let us know by creating an issue). Then add entries to [config.json](https://github.com/Johnnydouble/FileMango/blob/master/res/config.json) for each module interface. Filetypes are specefied as [mime types](https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types).

4. Autorun 

If you want to keep metadata constantly up to date, you must configure FileMango to start with your system. Either run the binary from a start script or set up a systemd service.

## Extended Attributes

In order to make use of the output of Filemango you must read the Extended Attributes of the file that the metadata is associated with. [Xattrvi](https://github.com/cherti/xattrvi) and the `getfattr` command (especially the `-R` and `-d` flags) are useful for debugging and simple use cases, most languages have a library to manage them for more complex situations.

## Authors

Created by Christen Spadavecchia and Oliver Easton
