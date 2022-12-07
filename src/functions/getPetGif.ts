const petpet = require('pet-pet-gif');


export async function getPetGif(content: string): Promise<Buffer | string> {
    let gif: any;
    try {
        gif = await petpet(content);
    } catch {
        gif = "ERROR";
    }
    return gif;
}
