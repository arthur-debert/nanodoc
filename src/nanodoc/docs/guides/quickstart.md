# Nanodoc

Nanodoc is a minimalist document bundler designed for stiching hints, reminders
and short docs. Useful for prompts, personalized docs highlights for your teams
or a note to your future self

No config, nothing to learn nor remember. Short , simple, sweet.

## Features

- No config, no tutorial, no pain.
- Combines multiple text files into a single document
- Adds clear title separators between pages
- Supports optional line numbering (per file or global)
- Can generate a table of contents
- Flexible file selection methods
- Customizable header styles and sequence numbering

## Usage

$ nanodoc file1.txt file2.txt

$ nanodoc -n file1.txt file2.txt # Per-file line numbering $ nanodoc -nn
file1.txt file2.txt # Global line numbering $ nanodoc -nn --toc file1.txt
file2.txt # Global numbering with TOC

## File Selection Options

Nanodoc is flexible in how you specify the files to bundle:

$ nanodoc `<file-1>....<file-n>` # individual files $ nanodoc `<dir-name>` # all
txt and md files in the dir will be included $ nanodoc `<dir-name> <file-1>` #
mix and match as you'd like $ nanodoc `<bundle>` # any .bundle.\* file that is a
list of paths, one per line $ nanodoc `<live-bundle>` # a file that mixes text
and file paths, where paths are replaced with their content

Get only parts of a file:

$ nanodoc readme.txt:L14-16,L30-50 # get the good parts only

## Command Line Options

- `-n`: Add per-file line numbering (01, 02, etc.)
- `-nn`: Add global line numbering: useful for referencing the full doc gen
  later
- `--toc`: Generate a table of contents at the beginning

## Get fancy

- `--seq`: numerical, roman or letter for ref the file sequence
- `--style`: nice (Human Readable (human-readable.txt), or file, or full-path

## Save for later

Generated a doc good enough to repeat, export the bundle

$nanodoc --export-bundle bestdocs.bundle.txt `<file-1>...<file-n>`

## Keep it simple

Nothing to config. Nothing to learn. No tutorials to watch.

In fact, you've just went through the full documentation. $ nanodoc --help # all
there is
