import { CommandInteraction } from "discord.js";
import { getPetGif } from "./getPetGif";
import { rateLimits } from "../index";

export async function createGif(interaction: CommandInteraction, content: string): Promise<Buffer | boolean> {
    let gif: Buffer | string;
    try {
        let options = {
            resolution: 128,
            delay: 30
        };
        if (interaction.options.getInteger("delay")) {
            options.delay = interaction.options.getInteger("delay")!;
        }
        if (interaction.options.getInteger("resolution")) {
            options.resolution = interaction.options.getInteger("resolution")!;
        }
        gif = await getPetGif(content);
        if (typeof gif === "string") { return false; }
        rateLimits.push({ id: interaction.user.id, time: Date.now() });
    } catch {
        interaction.editReply({
            content: "Sorry, but it looks like something went wrong. Please retry with a valid file/link",
        });
        return false;
    }
    return gif;
}
