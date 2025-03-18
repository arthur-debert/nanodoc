"""Command-line interface for nanodoc."""

import logging
import sys

import click

from . import VERSION
from .v1.core import run as run_v1

# Initialize logger
logger = logging.getLogger("nanodoc")

# Try importing v2 implementation
try:
    from .v2.core import run as run_v2

    V2_AVAILABLE = True
except ImportError:
    V2_AVAILABLE = False
    logger.warning("V2 implementation not available")

# Define Click context settings
CONTEXT_SETTINGS = {
    "help_option_names": ["-h", "--help"],
    "max_content_width": 100,
}


def setup_logging(verbose: bool) -> None:
    """Set up logging based on verbosity level."""
    if verbose:
        logger.setLevel(logging.DEBUG)
    else:
        logger.setLevel(logging.WARNING)


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
@click.option("--use-v2", is_flag=True, help="Use v2 implementation.")
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
    use_v2: bool,
) -> None:
    """Process source files and generate documentation."""
    setup_logging(verbose)

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
        # Choose implementation
        run_impl = run_v2 if use_v2 and V2_AVAILABLE else run_v1
        logger.info(f"Using {'v2' if use_v2 and V2_AVAILABLE else 'v1'} implementation")

        # Run the selected implementation with unified interface
        result = run_impl(
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
