"""Version information for nanodoc package."""

# Import version from package metadata
try:
    from importlib.metadata import version as get_version
except ImportError:
    # For Python < 3.8
    from importlib_metadata import version as get_version

# Get version from package metadata or fallback to pyproject.toml
try:
    VERSION = get_version("nanodoc")
except Exception:
    # Fallback to reading from pyproject.toml if package is not installed
    try:
        import pathlib

        import tomli

        # Try to find the pyproject.toml file
        file_path = pathlib.Path(__file__).parent.parent / "pyproject.toml"
        if file_path.exists():
            with open(file_path, "rb") as f:
                pyproject = tomli.load(f)
                VERSION = pyproject.get("project", {}).get("version", "unknown")
        else:
            VERSION = "unknown"
    except ImportError:
        VERSION = "unknown"


def get_version():
    """Return the version of nanodoc."""
    return VERSION
