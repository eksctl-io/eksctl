const Github_contributors_url = 'https://api.github.com/repos/eksctl-io/eksctl/contributors?per_page=100';

const getContributors = () => {
    axios
        .get(Github_contributors_url)
        .then((response) => {
            const contributors = response.data;
            const container = document.getElementById('contributors');
            const imageTags = contributors.map(img => `<a href="${img.html_url}" ><img src="${img.avatar_url}"></a>`);
            container.innerHTML = imageTags.join('');
        })
        .catch((error) => console.error(error));
};

document.addEventListener( 'DOMContentLoaded', function() {
    // Show Github Contributors on Homepage
    getContributors();

    // Show adopters carousel
    new Glide('.adopters', {
        type: 'carousel',
        autoplay: 4000,
        hoverpause: false,
        perView: 3,
        animationTimingFunc: 'linear',
        animationDuration: 1000
    }).mount()
});