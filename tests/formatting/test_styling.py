from nanodoc.formatting import apply_sequence_to_text, apply_style_to_filename


def test_apply_style_to_filename_with_filename_style():
    filename = "test_file.txt"
    original_path = "/path/to/test_file.txt"

    # Test with filename style
    result = apply_style_to_filename(filename, "filename", original_path)
    assert result == "test_file.txt"

    # Test with None style
    result = apply_style_to_filename(filename, None, original_path)
    assert result == "test_file.txt"


def test_apply_style_to_filename_with_path_style():
    filename = "test_file.txt"
    original_path = "/path/to/test_file.txt"

    result = apply_style_to_filename(filename, "path", original_path)
    assert result == "/path/to/test_file.txt"


def test_apply_style_to_filename_with_nice_style():
    filename = "test_file.txt"
    original_path = "/path/to/test_file.txt"

    result = apply_style_to_filename(filename, "nice", original_path)
    assert result == "Test File (test_file.txt)"

    # Test with hyphens and underscores
    filename = "test-file_name.txt"
    original_path = "/path/to/test-file_name.txt"

    result = apply_style_to_filename(filename, "nice", original_path)
    assert result == "Test File Name (test-file_name.txt)"


def test_apply_style_to_filename_without_original_path():
    filename = "test_file.txt"

    # Without original_path, should return filename regardless of style
    result = apply_style_to_filename(filename, "path", None)
    assert result == "test_file.txt"

    result = apply_style_to_filename(filename, "nice", None)
    assert result == "test_file.txt"


def test_apply_sequence_to_text_with_numerical_sequence():
    text = "Test Text"

    # Test with numerical sequence
    result = apply_sequence_to_text(text, "numerical", 0)
    assert result == "1. Test Text"

    result = apply_sequence_to_text(text, "numerical", 9)
    assert result == "10. Test Text"


def test_apply_sequence_to_text_with_letter_sequence():
    text = "Test Text"

    # Test with letter sequence
    result = apply_sequence_to_text(text, "letter", 0)
    assert result == "a. Test Text"

    result = apply_sequence_to_text(text, "letter", 1)
    assert result == "b. Test Text"

    # Test wrapping around after 'z'
    result = apply_sequence_to_text(text, "letter", 25)
    assert result == "z. Test Text"

    result = apply_sequence_to_text(text, "letter", 26)
    assert result == "a. Test Text"


def test_apply_sequence_to_text_with_roman_sequence():
    text = "Test Text"

    # Test with roman sequence
    result = apply_sequence_to_text(text, "roman", 0)
    assert result == "i. Test Text"

    result = apply_sequence_to_text(text, "roman", 1)
    assert result == "ii. Test Text"

    result = apply_sequence_to_text(text, "roman", 4)
    assert result == "v. Test Text"


def test_apply_sequence_to_text_with_none_sequence():
    text = "Test Text"

    # Test with None sequence
    result = apply_sequence_to_text(text, None, 0)
    assert result == "Test Text"

    # Test with invalid sequence
    result = apply_sequence_to_text(text, "invalid", 0)
    assert result == "Test Text"


def test_integration_of_style_and_sequence():
    filename = "test_file.txt"
    original_path = "/path/to/test_file.txt"

    # Apply style first, then sequence
    styled_text = apply_style_to_filename(filename, "nice", original_path)
    assert styled_text == "Test File (test_file.txt)"

    final_text = apply_sequence_to_text(styled_text, "numerical", 0)
    assert final_text == "1. Test File (test_file.txt)"

    # Different style and sequence
    styled_text = apply_style_to_filename(filename, "path", original_path)
    assert styled_text == "/path/to/test_file.txt"

    final_text = apply_sequence_to_text(styled_text, "roman", 2)
    assert final_text == "iii. /path/to/test_file.txt"
