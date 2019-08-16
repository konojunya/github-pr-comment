const Octokit = require("@octokit/rest");
const axios = require("axios");
const fs = require("fs");

const octokit = new Octokit({
  auth: process.env.GITHUB_ACCESS_TOKEN
});

async function main() {
  try {
    const json = JSON.parse(
      fs
        .readFileSync(process.env.GITHUB_EVENT_PATH, { encoding: "utf-8" })
        .toString()
    );
    const { pulls_url } = json.repository;
    const res = await axios.get(pulls_url, {
      headers: {
        Authorization: `token ${process.env.GITHUB_ACCESS_TOKEN}`
      }
    });

    console.log(res);
  } catch (e) {
    console.error(e);
    process.exit(1);
  }
}
main();

// async function createComment() {
//   const res = await octokit.issues.createComment({
//     owner: process.env.OWNER,
//     repo: process.env.REPO,
//     issue_number: await getIssueNumber(),
//     body: "commented from api! https://github.com"
//   });

//   console.log(res);
// }
