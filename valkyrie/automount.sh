mybucket=$1
mountPath=$2
passwdFile=$3
storageURL=$4

s3fs $mybucket $mountPath -o passwd_file=$passwdFile -o url=$storageURL -o use_path_request_style -o umask=0033 
#s3fs bucket1 /home/pine/fuse-mountpoint/bucket1 -o passwd_file=${HOME}/.passwd-s3fs -o url=http://localhost:9000 -o use_path_request_style -o dbglevel=info -f -o curldbg

