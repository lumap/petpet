import { CommandInteraction } from "discord.js";

export function sendGif(interaction: CommandInteraction, gif: Buffer, target: string) {
    interaction.editReply({
        files: [
            {
                attachment: gif,
                name: "pet.gif",
                description: `${interaction.user.tag} has pet ${target}`
            }
        ],
    }).then().catch(() => { });
}