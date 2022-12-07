import { CommandInteraction, MessageContextMenuInteraction, UserContextMenuInteraction } from "discord.js";

export function sendGif(interaction: CommandInteraction|UserContextMenuInteraction|MessageContextMenuInteraction, gif: Buffer) {
    interaction.editReply({
        files: [
            {
                attachment: gif,
                name: "pet.gif",
                description: "Pet!"
            }
        ],
    });
}
