"""Command-line interface for nanodoc using Click."""

import logging
import os
import pathlib
import sys

import click
from rich.console import Console

from .core import process_all
from .files import TXT_EXTENSIONS, get_files_from_args
from .formatting import create_themed_console, get_available_themes
from .version import VERSION

# Initialize console for rich output - will be updated with theme later
console = Console()

# Initialize logger
logger = logging.getLogger("nanodoc")
logger.setLevel(logging.INFO)  # Set default level to INFO


def setup_logging(to_stderr=False, enabled=False):
    """Configure logging based on requirements."""
    global logger
    if not logger.hasHandlers():
        level = logging.DEBUG if enabled else logging.INFO
        logger.setLevel(level)

        stream = sys.stderr if to_stderr else sys.stdout
        handler = logging.StreamHandler(stream)
        formatter = logging.Formatter(
            "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
        )
        handler.setFormatter(formatter)
        logger.addHandler(handler)
    else:
        level = logging.DEBUG if enabled else logging.INFO
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
    help="Line number mode (one -n for file, two for all)",
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
    help="Add file extensions to search for",
)
@click.option(
    "--theme",
    type=click.Choice(get_available_themes()),
    help="Select theme for rendering",
)
@click.argument("sources", nargs=-1, required=False)
@click.version_option(version=VERSION)
@click.pass_context
def cli(ctx, v, n, toc, no_header, sequence, style, txt_ext, theme, sources):
    """A minimalist document bundler for hints, reminders and docs."""
    # Set up logging based on verbose flag
    setup_logging(to_stderr=True, enabled=v)

    logger.info(f"CLI called with theme: {theme}")

    # Store theme in context for use by subcommands
    ctx.ensure_object(dict)
    ctx.obj["theme"] = theme

    # Update console with selected theme if provided
    global console
    if theme:
        logger.info(f"Creating console with theme: {theme}")
        console = create_themed_console(theme)

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
    help="Line number mode (one -n for file, two for all)",
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
    help="Add file extensions to search for",
)
@click.option(
    "--theme",
    type=click.Choice(get_available_themes()),
    help="Select theme for rendering",
)
@click.argument("sources", nargs=-1, required=True)
@click.pass_context
def process(ctx, v, n, toc, no_header, sequence, style, txt_ext, theme, sources):
    """Process source files and generate documentation.

    SOURCES are the files or directories to process.
    """
    try:
        # Set up logging based on verbose flag
        setup_logging(to_stderr=True, enabled=v)

        logger.info(f"Process command called with theme: {theme}")

        # Update console with theme if provided at command level
        global console
        if theme:
            logger.info(f"Creating console with theme from command: {theme}")
            console = create_themed_console(theme)
        # Otherwise use theme from parent context if available
        elif ctx.obj and "theme" in ctx.obj and ctx.obj["theme"]:
            theme_name = ctx.obj["theme"]
            logger.info(f"Creating console with theme from context: {theme_name}")
            console = create_themed_console(theme_name)

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
        logger.info("Rendering output with console")
        console.print(output)

    except Exception as e:
        err_msg = f"Error: {e}"
        logger.error(f"Error in process command: {e}", exc_info=True)
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
@click.option(
    "--theme",
    type=click.Choice(get_available_themes()),
    help="Select theme for rendering",
)
@click.pass_context
def show_help(ctx, topic, theme):
    """Show help for a specific topic.

    TOPIC is the name of the guide to display.
    """
    # Set up logging
    setup_logging(to_stderr=True)

    logger.info(f"Help command called with theme: {theme}")

    # Update console with theme if provided at command level
    global console
    if theme:
        logger.info(f"Creating console with theme from command: {theme}")
        console = create_themed_console(theme)
    # Otherwise use theme from parent context if available
    elif ctx.obj and "theme" in ctx.obj and ctx.obj["theme"]:
        theme_name = ctx.obj["theme"]
        logger.info(f"Creating console with theme from context: {theme_name}")
        console = create_themed_console(theme_name)

    if not topic:
        # Show general help using Rich for formatting
        from .help import _render_content, get_help_content

        found, content = get_help_content()
        if found:
            logger.info("Rendering help content with _render_content")
            _render_content(content)
        else:
            # Fallback to simple help if rich help content not found
            logger.info("Help content not found, using fallback")
            console.print("\n# nanodoc\n")
            msg = "A minimalist document bundler for hints, reminders and docs.\n"
            console.print(msg)

            # Import here to avoid circular imports
            from .help import get_options_section, get_topics_section

            # Show options section
            console.print("## OPTIONS\n")
            console.print(get_options_section())

            # Show topics section
            console.print("## HELP TOPICS\n")
            console.print(get_topics_section())
        return

    # Look for the topic in the guides directory
    from .help import get_guide_content

    logger.info(f"Looking for guide: {topic}")
    found, content = get_guide_content(topic)

    if found:
        logger.info(f"Guide found: {topic}")
        # For tests compatibility, use simple output format
        # This matches the expected format in the tests
        click.echo(f"\n{topic.upper()} GUIDE\n")
        click.echo("=" * 80)
        click.echo(content)
    else:
        logger.info(f"Guide not found: {topic}")
        # For tests compatibility, use simple error output
        click.echo(f"Guide '{topic}' not found. Available guides:", err=True)

        # Check if any guides exist
        guides_dir = _get_guides_dir()
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
@click.option(
    "--theme",
    type=click.Choice(get_available_themes()),
    help="Select theme for rendering",
)
@click.pass_context
def show_guide(ctx, topic, theme):
    """Show help for a specific topic (alias for 'help' command)."""
    # Set up logging
    setup_logging(to_stderr=True)

    logger.info(f"Guide command called with theme: {theme}")

    # Just call the help command with the same topic and theme
    ctx.invoke(show_help, topic=topic, theme=theme)


@cli.command("version")
def show_version():
    """Show the version of nanodoc."""
    click.echo(f"nanodoc version {VERSION}")


def main():
    """Main entry point for the nanodoc application."""
    try:
        # Set up basic logging
        setup_logging(to_stderr=True)

        # Check if the first argument is a command
        if len(sys.argv) > 1 and sys.argv[1] in cli.commands:
            command_name = sys.argv[1]
            logger.info(f"Running command: {command_name}")

            # Get the command
            command = cli.commands[command_name]

            # Run the command with the remaining arguments
            command.main(sys.argv[2:], standalone_mode=False)
        else:
            # Run the CLI normally
            logger.info("Running main CLI")
            cli(standalone_mode=False)
    except click.exceptions.Abort:
        # Handle keyboard interrupts gracefully
        logger.info("Command aborted")
        sys.exit(1)
    except click.exceptions.ClickException as e:
        # Handle Click exceptions
        logger.error(f"Click exception: {e}")
        e.show()
        sys.exit(e.exit_code)
    except Exception as e:
        # Handle other exceptions
        logger.error(f"Unexpected error: {e}", exc_info=True)
        click.echo(f"Error: {e}", err=True)
        sys.exit(1)


if __name__ == "__main__":
    main()
