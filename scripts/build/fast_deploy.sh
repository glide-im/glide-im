mkdir "/srv/glideim"

cd /srv/glideim || exit
wget https://github.91chi.fun//https://github.com/dengzii/go_im/releases/latest/download/go_im_singleton_$1.tar.gz -O go_im_$1.tar.gz
tar zxvf go_im_$1.tar.gz
pkill go_im
./go_im
rm go_im_$1.tar.gz