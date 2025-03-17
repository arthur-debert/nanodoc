import os
import sys

import pytest

# Add parent directory to path to import nanodoc
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from nanodoc.formatting import to_roman


def test_basic_roman_conversions():
    """Test basic roman numeral conversions."""
    test_cases = [
        (1, "i"),
        (2, "ii"),
        (3, "iii"),
        (4, "iv"),
        (5, "v"),
        (9, "ix"),
        (10, "x"),
        (14, "xiv"),
        (19, "xix"),
        (24, "xxiv"),
        (42, "xlii"),
        (99, "xcix"),
        (501, "di"),
        (1984, "mcmlxxxiv"),
    ]

    for num, expected in test_cases:
        assert to_roman(num) == expected, f"Failed converting {num} to roman"


def test_zero_input():
    """Test that zero input raises ValueError."""
    with pytest.raises(ValueError):
        to_roman(0)


def test_negative_input():
    """Test that negative input raises ValueError."""
    with pytest.raises(ValueError):
        to_roman(-1)


def test_non_integer_input():
    """Test that non-integer inputs raise ValueError."""
    with pytest.raises(ValueError):
        to_roman("not a number")

    with pytest.raises(ValueError):
        to_roman(3.14)
