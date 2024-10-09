from io import BytesIO
from petpetgif import petpet

def make_petpet_gif(bytes: BytesIO, resolution: int = 128) -> bytes:
    petpet.resolution = (resolution, resolution)
    
    output_bytes = BytesIO()

    petpet.make(bytes, output_bytes)

    petpet.resolution = (128,128) # It doesn't reset when a gif is done for some reason

    return output_bytes.getvalue()