const Octokit = require("@octokit/rest");
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
    const splited = json.repository.pulls_url.split(/\//);
    const pull_request_number = splited[splited.length - 1];
    const comment = await getComment(pull_request_number);
    if (comment == null) {
      createComment(pull_request_number);
    } else {
      updateComment(comment.id);
    }

    console.log("done");
  } catch (e) {
    console.error(e);
    process.exit(1);
  }
}
main();

async function createComment(issue_number) {
  const res = await octokit.issues.createComment({
    owner: process.env.OWNER,
    repo: process.env.REPO,
    issue_number,
    body: "commented from api! https://github.com"
  });

  console.log(res);
}

async function updateComment(comment_id) {
  octokit.issues.updateComment({
    owner: process.env.OWNER,
    repo: process.env.REPO,
    comment_id,
    body: `updated!: ${new Date()}`
  });
}

async function getComment(issue_number) {
  const { data } = await octokit.issues.listComments({
    owner: process.env.OWNER,
    repo: process.env.REPO,
    issue_number
  });

  return data.find(comment => comment.user.login === "konojunya");
}
