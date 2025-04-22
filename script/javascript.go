package script

import "fmt"

const (
	VoteButtonSelector  = "button[aria-label^=\"赞同 \"]"
	VotedButtonSelector = "button[aria-label^=\"已赞同 \"]"
	LikeButtonSelector  = "button[aria-live^=\"polite\"]"
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

func GetLikeScript() string {
	return fmt.Sprintf(`
		(() => {
			// execute script
			const buttons = Array.from(document.querySelectorAll('%s'));
			if (buttons.length === 0) {
				return false;
			}

			const targetBtn = buttons.find(btn => btn.textContent.includes('喜欢') && !btn.textContent.includes('取消喜欢'));
			if (!targetBtn) {
				return false;
			}

			targetBtn.click();
		})();
	`, LikeButtonSelector)
}

func GetCheckIfLikedScript() string {
	return fmt.Sprintf(`
		(() => {
			// check if the button is clicked
			const buttons = Array.from(document.querySelectorAll('%s'));
			const button = buttons.find(btn => btn.textContent.includes('取消喜欢'));

			if (button) {
				return true;
			}
			else {
				return false;
			}
		})();
	`, LikeButtonSelector)
}
