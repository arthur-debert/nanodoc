"""Command-line interface for nanodoc using Click."""

import logging
import os
import pathlib
import sys

import click
from rich.console import Console

from .core import process_all
from .files import TXT_EXTENSIONS, get_files_from_args
from .version import VERSION

# Initialize console for rich output
console = Console()

# Initialize logger
logger = logging.getLogger("nanodoc")
logger.setLevel(logging.CRITICAL)  # Start with logging disabled


def setup_logging(to_stderr=False, enabled=False):
    """Configure logging based on requirements."""
    global logger
    if not logger.hasHandlers():
        level = logging.DEBUG if enabled else logging.CRITICAL
        logger.setLevel(level)

        stream = sys.stderr if to_stderr else sys.stdout
        handler = logging.StreamHandler(stream)
        formatter = logging.Formatter(
            "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
        )
        handler.setFormatter(formatter)
        logger.addHandler(handler)
    else:
        level = logging.DEBUG if enabled else logging.CRITICAL
        logger.setLevel(level)
    return logger


# Define Click context settings
CONTEXT_SETTINGS = {
    "help_option_names": ["-h", "--help"],
    "max_content_width": 100,
}


# Custom Group class that handles the "help" command differently
class NanodocGroup(click.Group):
    def get_command(self, ctx, cmd_name):
        # First try to get the command normally
        rv = click.Group.get_command(self, ctx, cmd_name)
        if rv is not None:
            return rv

        # If the command is not found and it's "help", handle it specially
        if cmd_name == "help":
            return self.commands["help"]

        return None


# For backward compatibility, we need to handle the case where options are passed
# directly to the main command instead of to the 'process' subcommand
@click.group(
    cls=NanodocGroup, context_settings=CONTEXT_SETTINGS, invoke_without_command=True
)
@click.option("-v", is_flag=True, help="Enable verbose mode")
@click.option(
    "-n",
    count=True,
    help="Enable line number mode (one -n for file, two for all)",
)
@click.option("--toc", is_flag=True, help="Generate table of contents")
@click.option("--no-header", is_flag=True, help="Hide file headers")
@click.option(
    "--sequence",
    type=click.Choice(["numerical", "letter", "roman"]),
    help="Add sequence numbers to headers",
)
@click.option(
    "--style",
    type=click.Choice(["filename", "path", "nice"]),
    default="nice",
    help="Header style: nice, filename or path",
)
@click.option(
    "--txt-ext",
    multiple=True,
    help="Add additional file extensions to search for",
)
@click.argument("sources", nargs=-1, required=False)
@click.version_option(version=VERSION)
@click.pass_context
def cli(ctx, v, n, toc, no_header, sequence, style, txt_ext, sources):
    """A minimalist document bundler for hints, reminders and docs."""
    # If no subcommand is provided, run the process command with the given options
    if ctx.invoked_subcommand is None and sources:
        # Call the process command with the provided options
        ctx.invoke(
            process,
            v=v,
            n=n,
            toc=toc,
            no_header=no_header,
            sequence=sequence,
            style=style,
            txt_ext=txt_ext,
            sources=sources,
        )
    elif ctx.invoked_subcommand is None:
        # If no sources and no subcommand, show help
        click.echo(ctx.get_help())


@cli.command()
@click.option("-v", is_flag=True, help="Enable verbose mode")
@click.option(
    "-n",
    count=True,
    help="Enable line number mode (one -n for file, two for all)",
)
@click.option("--toc", is_flag=True, help="Generate table of contents")
@click.option("--no-header", is_flag=True, help="Hide file headers")
@click.option(
    "--sequence",
    type=click.Choice(["numerical", "letter", "roman"]),
    help="Add sequence numbers to headers",
)
@click.option(
    "--style",
    type=click.Choice(["filename", "path", "nice"]),
    default="nice",
    help="Header style: nice, filename or path",
)
@click.option(
    "--txt-ext",
    multiple=True,
    help="Add additional file extensions to search for",
)
@click.argument("sources", nargs=-1, required=True)
def process(v, n, toc, no_header, sequence, style, txt_ext, sources):
    """Process source files and generate documentation.

    SOURCES are the files or directories to process.
    """
    try:
        # Set up logging based on verbose flag
        setup_logging(to_stderr=True, enabled=v)

        # Process line numbering mode
        if n == 0:
            line_number_mode = None
        elif n == 1:
            line_number_mode = "file"
        else:  # n >= 2
            line_number_mode = "all"

        # Process additional file extensions if provided
        extensions = list(TXT_EXTENSIONS)  # Create a copy of the default extensions
        if txt_ext:
            for ext in txt_ext:
                # Add a leading dot if not present
                if not ext.startswith("."):
                    ext = "." + ext
                # Only add if not already in the list
                if ext not in extensions:
                    extensions.append(ext)

        # Get verified content items from arguments
        if txt_ext:
            # Only pass extensions if custom extensions were provided
            content_items = get_files_from_args(sources, extensions=extensions)
        else:
            # Use default extensions
            content_items = get_files_from_args(sources)

        # Process the files and print the result
        if not content_items:
            click.echo("Error: No valid source files found.", err=True)
            sys.exit(1)

        output = process_all(
            content_items,
            line_number_mode,
            toc,
            not no_header,
            sequence,
            style,
        )
        click.echo(output)

    except Exception as e:
        err_msg = f"Error: {e}"
        click.echo(err_msg, err=True)
        sys.exit(1)


def _get_docs_dir():
    """Return the path to the docs directory."""
    module_dir = pathlib.Path(__file__).parent.absolute()
    return module_dir / "docs"


def _get_guides_dir():
    """Return the path to the guides directory."""
    module_dir = pathlib.Path(__file__).parent.absolute()
    guides_dir = module_dir / "docs" / "guides"

    # Create the guides directory if it doesn't exist
    if not guides_dir.exists():
        os.makedirs(guides_dir, exist_ok=True)

    return guides_dir


@cli.command("help")
@click.argument("topic", required=False)
def show_help(topic):
    """Show help for a specific topic.

    TOPIC is the name of the guide to display.
    """
    if not topic:
        # Show general help
        click.echo("\n# nanodoc\n")
        click.echo("A minimalist document bundler for hints, reminders and docs.\n")

        # Import here to avoid circular imports
        from .help import get_options_section, get_topics_section

        # Show options section
        click.echo("## OPTIONS\n")
        click.echo(get_options_section())

        # Show topics section
        click.echo("## HELP TOPICS\n")
        click.echo(get_topics_section())
        return

    # Look for the topic in the guides directory
    guides_dir = _get_guides_dir()
    found = False

    # Check for the guide with various extensions
    for ext in TXT_EXTENSIONS:
        guide_path = guides_dir / f"{topic}{ext}"
        if guide_path.exists():
            found = True
            with open(guide_path, "r", encoding="utf-8") as f:
                content = f.read()

            # Use Click's formatting for simple output
            click.echo(f"\n{topic.upper()} GUIDE\n")
            click.echo("=" * 80)
            click.echo(content)
            break

    if not found:
        # List available guides
        click.echo(
            f"Guide '{topic}' not found. Available guides:",
            err=True,
        )

        # Check if any guides exist
        guides_exist = False
        for ext in TXT_EXTENSIONS:
            for guide_path in guides_dir.glob(f"*{ext}"):
                guides_exist = True
                guide_name = guide_path.name.replace(ext, "")
                click.echo(f"  - {guide_name}", err=True)

        if not guides_exist:
            click.echo("  No guides available.", err=True)

        sys.exit(1)


@cli.command("guide")
@click.argument("topic", required=False)
def show_guide(topic):
    """Show help for a specific topic (alias for 'help' command)."""
    # Just call the help command with the same topic
    ctx = click.get_current_context()
    ctx.invoke(show_help, topic=topic)


@cli.command("version")
def show_version():
    """Show the version of nanodoc."""
    click.echo(f"nanodoc version {VERSION}")


def main():
    """Main entry point for the nanodoc application."""
    try:
        # Check if the first argument is a command
        if len(sys.argv) > 1 and sys.argv[1] in cli.commands:
            command_name = sys.argv[1]

            # Get the command
            command = cli.commands[command_name]

            # Run the command with the remaining arguments
            command.main(sys.argv[2:], standalone_mode=False)
        else:
            # Run the CLI normally
            cli(standalone_mode=False)
    except click.exceptions.Abort:
        # Handle keyboard interrupts gracefully
        sys.exit(1)
    except click.exceptions.ClickException as e:
        # Handle Click exceptions
        e.show()
        sys.exit(e.exit_code)
    except Exception as e:
        # Handle other exceptions
        click.echo(f"Error: {e}", err=True)
        sys.exit(1)


if __name__ == "__main__":
    main()
