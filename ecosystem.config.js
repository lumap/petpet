module.exports = {
    apps: [
        {
            name: "petpet-prod",
            script: "npx tsc && npm run prod"
        },
        {
            name: "petpet-dev",
            script: "npx tsc && npm run dev"
        }
    ]
}