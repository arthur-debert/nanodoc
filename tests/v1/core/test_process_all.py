from nanodoc.v1.core import run_all
from nanodoc.v1.data import ContentItem, LineRange


def test_process_all_toc_generation():
    """Test TOC generation with multiple content items."""
    content_items = [
        ContentItem(
            file_path="test1.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 1\n", "Line 2\n"],
        ),
        ContentItem(
            file_path="test2.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 3\n", "Line 4\n"],
        ),
    ]

    result = run_all(
        content_items,
        line_number_mode=None,
        generate_toc=True,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    assert "Table of Contents" in result
    assert "test1.txt" in result
    assert "test2.txt" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result


def test_process_all_with_no_header():
    """Test processing multiple content items without headers."""
    content_items = [
        ContentItem(
            file_path="test1.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 1\n", "Line 2\n"],
        ),
        ContentItem(
            file_path="test2.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 3\n", "Line 4\n"],
        ),
    ]

    result = run_all(
        content_items,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=False,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    assert "test1.txt" not in result
    assert "test2.txt" not in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result


def test_process_all_with_header_sequence():
    """Test processing multiple content items with header sequence."""
    content_items = [
        ContentItem(
            file_path="test1.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 1\n", "Line 2\n"],
        ),
        ContentItem(
            file_path="test2.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 3\n", "Line 4\n"],
        ),
    ]

    result = run_all(
        content_items,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence="numerical",
        style=None,
        txt_ext=None,
    )

    assert "1. test1.txt" in result
    assert "2. test2.txt" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result


def test_process_all_with_header_style():
    """Test processing multiple content items with header style."""
    content_items = [
        ContentItem(
            file_path="test1.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 1\n", "Line 2\n"],
        ),
        ContentItem(
            file_path="test2.txt",
            ranges=[LineRange(start=1, end=2)],
            content=["Line 3\n", "Line 4\n"],
        ),
    ]

    result = run_all(
        content_items,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style="nice",
        txt_ext=None,
    )

    assert "Test1 (test1.txt)" in result
    assert "Test2 (test2.txt)" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result
