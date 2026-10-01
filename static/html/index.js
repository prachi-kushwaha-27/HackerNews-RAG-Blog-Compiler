function getAvailablePostDates(next) {
    fetch('/available-post-dates')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(dates => {
            next(dates);
        })
        .catch(error => {
            console.error('error fetching post data:', error);
        });
}

function getSummary(postId, date, next) {
    fetch('/summary?id=' + postId + '&date=' + date)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(summary => {
            next(summary);
        })
        .catch(error => {
            console.error('error fetching post data:', error);
        });
}

function getPosts(date, next) {
    fetch('/posts?date=' + date)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(posts => {
            next(posts);
        })
        .catch(error => {
            console.error('error fetching post data:', error);
        });
}

function dateFromDatetime(date) {
    return date.substring(0, 10)
}

function formatDateForDisplay(date) {
    var options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
    return new Date(date).toLocaleString('en-US', options)
}
