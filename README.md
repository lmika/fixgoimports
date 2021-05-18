# fixgoimports

A very simple go importing formatter.  This sorts and groups imports alphabetically with standard library
imports appearing above 3rd party imports.

## Usage

```
fixgoimports [-w] [FILE_OR_DIR...]
```

Where:

- `-w`: writes the formatted source file back to the original file

When called with arguments, each argument must either be a file or a directory.  Files will be processed as long
as they have the suffix `.go`.  Directories will be traversed, including any subdirectories, minus any file or
directory that begins with `.` or `_`, or any directory with the name `vendor`.

## Limitations

- Only recognises `import ( ... )` blocks and can only handle a single import block at a time.
- Cannot determine whether an import block appears in a string or not.
