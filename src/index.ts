import { Client, Intents } from "discord.js";
import { handleMessageContextMenu } from "./functions/handleMessageContextMenu";
import { handleUserContextMenu } from "./functions/handleUserContextMenu";
import { handleSlashCommand } from "./functions/handleSlashCommand";
import { isRatelimited } from "./functions/isRatelimited";
import { config } from "../config"


const client = new Client({ intents: [Intents.FLAGS.GUILDS] });

export let rateLimits: { time: number, id: string }[] = [];

setInterval(function () {
    rateLimits = rateLimits.filter(c => c.time + 60000 > Date.now())
}, 60000)

client.on("interactionCreate", async (interaction) => {
    if (isRatelimited(interaction)) return;
    if (interaction.isUserContextMenu()) {
        handleUserContextMenu(interaction, client);
    } else if (interaction.isMessageContextMenu()) {
        handleMessageContextMenu(interaction, client)
    } else { //slash command
        if (!interaction.isCommand()) return;
        await handleSlashCommand(interaction)
    }
})

client.login(config.token);

client.on("ready", () => {
    console.log("bot started ig")
})