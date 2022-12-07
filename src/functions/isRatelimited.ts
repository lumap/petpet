import { Interaction } from "discord.js";
import { rateLimits } from "../index";

export function isRatelimited(interaction: Interaction): boolean {
    if (!interaction.isCommand() && !interaction.isUserContextMenu() && !interaction.isMessageContextMenu()) {
        console.log("L")
        return true;
    }
    const userRatelimits = rateLimits.filter(c => c.time + 60000 > Date.now()).filter(c => c.id === interaction.user.id);
    if (userRatelimits.length > 5) {
        interaction.reply({
            ephemeral: true,
            content: "You've been ratelimited (1 petpet per 10 seconds). This system has been put in place to prevent spam and power consuption from my server. Wait ~1 minute before being able to use the bot again."
        });
        return true;
    }
    return false;
}
