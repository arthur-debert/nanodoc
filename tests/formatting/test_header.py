from nanodoc.formatting import create_header


def test_create_header_contains_text():
    text = "Test Header"
    header = create_header(text)
    assert text in header
    assert header == text


def test_create_header_with_custom_char():
    text = "Test Header"
    char = "*"
    header = create_header(text, char)
    # Custom char is no longer used in the header
    assert header == text


def test_create_header_empty_text():
    text = ""
    header = create_header(text)
    assert header == text


def test_create_header_with_numerical_sequence():
    text = "Test Header"
    header = create_header(text, sequence="numerical", seq_index=0)
    assert header == "1. Test Header"

    header = create_header(text, sequence="numerical", seq_index=1)
    assert header == "2. Test Header"


def test_create_header_with_letter_sequence():
    text = "Test Header"
    header = create_header(text, sequence="letter", seq_index=0)
    assert header == "a. Test Header"

    header = create_header(text, sequence="letter", seq_index=1)
    assert header == "b. Test Header"


def test_create_header_with_roman_sequence():
    text = "Test Header"
    header = create_header(text, sequence="roman", seq_index=0)
    assert header == "i. Test Header"

    header = create_header(text, sequence="roman", seq_index=1)
    assert header == "ii. Test Header"


def test_create_header_with_filename_style():
    text = "test_file.txt"
    header = create_header(
        text, style="filename", original_path="/path/to/test_file.txt"
    )
    assert header == "test_file.txt"


def test_create_header_with_path_style():
    text = "test_file.txt"
    header = create_header(text, style="path", original_path="/path/to/test_file.txt")
    assert header == "/path/to/test_file.txt"


def test_create_header_with_nice_style():
    text = "test_file.txt"
    header = create_header(text, style="nice", original_path="/path/to/test_file.txt")
    assert header == "Test File (test_file.txt)"


def test_create_header_long_text():
    text = "This is a very long header that exceeds the maximum line width"
    header = create_header(text)
    assert header == text
