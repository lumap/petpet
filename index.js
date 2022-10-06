const { token } = require("./token.json")
const Discord = require('discord.js');
const petpet = require('pet-pet-gif');
let urlcheck = require('is-a-url');
const { parse } = require('twemoji-parser')

const client = new Discord.Client({ intents: [Discord.Intents.FLAGS.GUILDS] });

function isImage(url) {
    return /\.(jpg|jpeg|png)$/.test(url.toLowerCase());
}

let rateLimits = [];

client.on("interactionCreate", async (interaction, _member) => {
    if (!interaction.isCommand()) return;
    rateLimits = rateLimits.filter(c => c.time + 60000 > Date.now());
    const userRatelimits = rateLimits.filter(c => c.id === interaction.user.id);
    if (userRatelimits.length > 5) {
        return interaction.reply({
            ephemeral: true,
            content: "You've been ratelimited (1 petpet per 10 seconds). This system has been put in place to prevent spam and power consuption from my server. Wait ~1 minute before being able to use the bot again."
        });
    }
    try {
        await interaction.deferReply();
    } catch {
        return;
    }
    let content;
    if (interaction.options.getSubcommand() == "user") {
        if (interaction.guild) {
            try {
                const member = await interaction.guild.members.fetch(interaction.options.getUser("user"));
                content = member.displayAvatarURL({ format: "png", size: 1024 })
            } catch {
                content = "This member doesn't seem to be here. If you want to petpet them, do it in my DMs.";
            }
        } else {
            content = interaction.options.getUser("user").avatarURL({ format: "png", size: 1024 });
        }
    } else if (interaction.options.getSubcommand() == "attachment") {
        const url = interaction.options.getAttachment("attachment").url;
        if (!urlcheck(url) || !isImage(url)) {
            content = "Sorry, this attachment does not seem to be valid. Please make sure it's a `jpg`, `jpeg` or `png` image."
        } else {
            content = url;
        }
    } else if (interaction.options.getSubcommand() == "imageurl") {
        const url = interaction.options.getString("url");
        if (!urlcheck(url) || url.startsWith("https://tenor.com/view/") || !isImage(url)) {
            content = "Sorry, this link does not seem to be valid. Please make sure the image link ends with `.jpg`, `.jpeg` or `.png`."
        } else {
            content = url;
        }
    } else if (interaction.options.getSubcommand() == "emoji") {
        const emoji = interaction.options.getString("emoji");
        if ((emoji.match(/(<a?)?:\w+:(\d{16,20}>)?/u)) !== null) {
            content = `https://cdn.discordapp.com/emojis/${emoji.split(":")[2].slice(0, -1)}.png`
        } else if (parse(emoji)?.[0]?.url) {
            content = parse(emoji)[0].url
        } else {
            content = "I wasn't able to find an emoji in this. I wish discord had an \"emoji\" option for slash commands"
        }
    } else {
        return interaction.editReply({
            content: "how did you fuck that up"
        })
    }
    if (!content.startsWith("http")) {
        return interaction.editReply({
            content: content
        })
    }
    let gif;
    try {
        gif = await petpet(content);
        rateLimits.push({ id: interaction.user.id, time: Date.now() });
    } catch {
        return interaction.editReply({
            content: "Sorry, but it looks like something went wrong. Please retry with a valid file/link",
        })
    }
    interaction.editReply({
        files: [
            {
                attachment: gif,
                name: "pet.gif",
                description: "Pet!"
            }
        ]
    });
})

client.login(token);

client.on("ready", () => {
    console.log("bot started ig")
})