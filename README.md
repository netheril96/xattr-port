# Overview

A simple command line tool to export/import all extended attributes under a directory.

# Usage

The tool requires a directory path and a database file path for its operations.

## Exporting Extended Attributes

To export all extended attributes from files under a specific directory into a SQLite database:

```bash
xattr-port export --dir /path/to/your/directory --db /path/to/your/xattrs.db
```

## Importing Extended Attributes

To import extended attributes from a SQLite database and apply them to files under a specific directory:

```bash
xattr-port import --dir /path/to/your/directory --db /path/to/your/xattrs.db
