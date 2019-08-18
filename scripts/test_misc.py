import api


def test_misc():
    assert api.is_uid('0x10') == True
