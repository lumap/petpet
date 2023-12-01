const petpet = require('pet-pet-gif');

export async function makePetGif(content: string, options: any): Promise<Buffer | string> {
    let gif: any;
    try {
        gif = await petpet(content, options);
    } catch (e) {
        gif = "shit";
    }
    return gif;
}
