"""Command-line interface for nanodoc."""

import logging
import sys

import click

from . import VERSION
from .boot import configure_logging
from .core import run

# Initialize logger
logger = logging.getLogger("nanodoc")

# Define Click context settings
CONTEXT_SETTINGS = {
    "help_option_names": ["-h", "--help"],
    "max_content_width": 100,
}


@click.command()
@click.argument("sources", nargs=-1, type=click.Path(exists=True))
@click.option("-v", "--verbose", is_flag=True, help="Enable verbose output.")
@click.option("--toc", is_flag=True, help="Generate table of contents.")
@click.option("-n", count=True, help="Line number mode (one -n for file, two for all)")
@click.option("--theme", type=str, help="Theme to use for output.")
@click.option("--no-header", is_flag=True, help="Don't show file headers.")
@click.option("--sequence", type=click.Choice(["numerical", "letter", "roman"]))
@click.option(
    "--style",
    type=click.Choice(["filename", "path", "nice"]),
    default="nice",
    help="Header style",
)
@click.option("--txt-ext", multiple=True, help="Add file extensions to search for")
@click.version_option(version=VERSION)
def main(
    sources: list[str],
    verbose: bool,
    toc: bool,
    n: int,
    theme: str,
    no_header: bool,
    sequence: str,
    style: str,
    txt_ext: list[str],
) -> None:
    """Process source files and generate documentation."""
    configure_logging(verbose)

    if not sources:
        click.echo("No source files provided.", err=True)
        sys.exit(1)

    # Convert -n/-nn to line number mode
    line_number_mode = None
    if n == 1:
        line_number_mode = "file"
    elif n >= 2:
        line_number_mode = "all"

    try:
        # Using v2 implementation
        logger.info("Using v2 implementation")

        # Run implementation with unified interface
        result = run(
            sources=list(sources),
            line_number_mode=line_number_mode,
            generate_toc=toc,
            theme=theme,
            show_header=not no_header,
            sequence=sequence,
            style=style,
            txt_ext=txt_ext,
        )

        click.echo(result)

    except Exception as e:
        logger.error(str(e))
        sys.exit(1)


if __name__ == "__main__":
    main()
