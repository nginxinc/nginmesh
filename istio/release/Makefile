VERSION=0.2.12
REPO=nginmesh
USER=nginmesh
RELEASE_FILE=nginmesh-${VERSION}


update-install:
	./updateVersion.sh


release:	release_tar create_release upload_release

clean_release:  
	rm -rf ${RELEASE_FILE} ${RELEASE_FILE}.tar.gz

release_tar:    clean_release
	mkdir -p ${RELEASE_FILE}
	cp -r samples ${RELEASE_FILE}
	cp -r install ${RELEASE_FILE}
	cp -r ../tools ${RELEASE_FILE}
	rm -rf ${RELEASE_FILE}/install/kubernetes/templates
	rm -rf ${RELEASE_FILE}/install/kubernetes/.gitignore
	rm -rf ${RELEASE_FILE}/install/kubernetes/*.bak
	cd target;tar -czvf ${RELEASE_FILE}.tar.gz ${RELEASE_FILE}

create_release:
	github-release release \
    --user ${USER} \
    --repo ${REPO} \
    --tag ${VERSION} \
    --name "${VERSION}"

upload_release:
	github-release upload \
    	--user ${USER} \
    	--repo ${REPO} \
    	--tag ${VERSION} \
    	--name "${RELEASE_FILE}.tar.gz" \
    	--file ${RELEASE_FILE}.tar.gz



delete_release:
	github-release delete \
    --user ${USER} \
    --repo ${REPO} \
    --tag ${VERSION}



clean:	clean_release
