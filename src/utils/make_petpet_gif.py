from io import BytesIO
from src.pet_gen import petpet

def convert_speed_to_petpet_delay(speed: int) -> int:
    return int(40 // speed)

def make_petpet_gif(bytes: BytesIO, resolution: int = 128, speed = 1) -> bytes:
    petpet.resolution = (resolution, resolution)
    petpet.delay = convert_speed_to_petpet_delay(speed)
    
    output_bytes = BytesIO()

    petpet.make(bytes, output_bytes)

    # set values back to default
    petpet.resolution = (128,128)
    petpet.delay = 40

    return output_bytes.getvalue()