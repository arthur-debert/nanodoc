"""Utilities for test fixtures and common test operations."""

from pathlib import Path

from nanodoc.structures import FileContent


def get_fixtures_dir() -> Path:
    """Get the path to the fixtures directory."""
    return Path(__file__).parent / "fixtures"


def get_fixture_path(fixture_name: str) -> Path:
    """Get the full path to a fixture file.

    Args:
        fixture_name: Name of the fixture file (e.g. "cake.txt")

    Returns:
        Path to the fixture file
    """
    fixture_path = get_fixtures_dir() / fixture_name
    if not fixture_path.exists():
        raise FileNotFoundError(
            f"Fixture file {fixture_name} not found at {fixture_path}"
        )
    return fixture_path


def read_fixture(fixture_name: str) -> str:
    """Read the contents of a fixture file.

    Args:
        fixture_name: Name of the fixture file (e.g. "cake.txt")

    Returns:
        Contents of the fixture file as a string
    """
    with open(get_fixture_path(fixture_name)) as f:
        return f.read()


def create_fixture_content_item(fixture_name: str) -> FileContent:
    """Create a FileContent from a fixture file.

    This is the preferred way to access fixture files in tests, as it provides
    a complete FileContent object that can be used with the nanodoc functions.

    Args:
        fixture_name: Name of the fixture file (e.g. "cake.txt")

    Returns:
        FileContent instance representing the fixture file
    """
    path = get_fixture_path(fixture_name)
    with open(path) as f:
        content = f.read()

    # Calculate the number of lines
    num_lines = len(content.splitlines()) or 1

    file_content = FileContent(
        filepath=str(path),
        ranges=[(1, num_lines)],
        content=content,
        is_bundle=fixture_name.endswith(".ndoc"),
        original_source=fixture_name,
    )

    # Add file_path property for compatibility with v1 tests
    file_content.file_path = file_content.filepath

    return file_content
