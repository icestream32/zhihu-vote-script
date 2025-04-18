package script

import "fmt"

const (
	VoteButtonSelector  = "button[aria-label^=\"赞同 \"]"
	VotedButtonSelector = "button[aria-label^=\"已赞同 \"]"
	LoginScript         = "document.querySelector('#root div.Popover.AppHeader-menu') !== null"
)

func GetVoteScript() string {
	return fmt.Sprintf(`
		(() => {

			const voteButtons = Array.from(document.querySelectorAll('%s'));
			if (voteButtons.length === 0) {
				throw new Error("No vote button found");
			}

			const validButton = voteButtons.find(button => !button.classList.contains('is-active'));
			if (!validButton) {
				throw new Error("No clickable vote button found");
			}

			validButton.click();

			const votedbuttons = Array.from(document.querySelectorAll('%s'));
			if (votedbuttons.length === 0) {
				throw new Error("Vote failed, no active vote button found");
			}

			return votedbuttons.some(button => button.classList.contains('is-active'));
		})();
	`, VoteButtonSelector, VotedButtonSelector)
}
