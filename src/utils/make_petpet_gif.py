from io import BytesIO
from petpetgif import petpet
from utils.url_to_bytesio import file_url_to_bytesio

def make_petpet_gif(url: str, resolution: int = 128) -> bytes:
    petpet.resolution = (resolution, resolution)

    avatar_bytes = file_url_to_bytesio(url)
    output_bytes = BytesIO()

    petpet.make(avatar_bytes, output_bytes)

    petpet.resolution = (128,128) # It doesn't reset when a gif is done for some reason

    return output_bytes.getvalue()