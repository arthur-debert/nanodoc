# nanodoc

An ultra-lightweight documentation generator that combines text files into a
single document.

USAGE nanodoc [options] [file1.txt file2.txt ...]

FORMATTING --toc: Include a table of contents at the beginning --no-header: Hide
file headers completely --sequence: Add sequence numbers to headers (numerical,
letter, or roman) --style: Change how filenames are displayed (filename, path,
nice)

LINE NUMBERS -n: Enable per-file line numbering -nn: Enable global line
numbering

ADDITIONAL --txt-ext: Add additional file extensions to search for

HELP TOPICS manifesto: Nanodoc Manifesto quickstart: Nanodoc Quick Start Guide

FLAGS -v, --verbose Enable verbose output --help Show help for command --version
Show nanodoc version

EXAMPLES $ nanodoc file1.txt file2.txt $ nanodoc -n file1.txt file2.txt #
Per-file line numbering $ nanodoc -nn file1.txt file2.txt # Global line
numbering $ nanodoc -nn --toc file1.txt file2.txt # Global numbering with TOC $
nanodoc dir-name # All txt and md files in the dir will be included $ nanodoc
dir-name file-1 # Mix and match as you'd like $ nanodoc bundle # Any .bundle.\*
file that is a list of paths, one per line $ nanodoc readme.txt:L14-16,L30-50 #
Get the good parts only

CORE COMMANDS help: Show help information or specific guides version: Display
the current version of nanodoc

LEARN MORE Use `nanodoc help <guide-name>` for more information about a specific
topic.
