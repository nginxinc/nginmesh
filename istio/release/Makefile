VERSION=0.6.0-RC1
REPO=nginmesh
USER=nginmesh
RELEASE_FILE=nginmesh-${VERSION}
RELEASE_DIR=target

release:	release_tar create_release upload_release

clean_release: 
	rm -rf ${RELEASE_DIR}

release_tar:    clean_release
	mkdir -p ${RELEASE_DIR}/${RELEASE_FILE}
	cp -r samples ${RELEASE_DIR}/${RELEASE_FILE}
	cp -r install ${RELEASE_DIR}/${RELEASE_FILE}
	cp -r tools ${RELEASE_DIR}/${RELEASE_FILE}
	rm -rf ${RELEASE_DIR}/${RELEASE_FILE}/install/kubernetes/.gitignore
	rm -rf ${RELEASE_DIR}/${RELEASE_FILE}/install/kubernetes/generated
	rm -rf ${RELEASE_DIR}/${RELEASE_FILE}/install/kubernetes/*.bak
	cd ${RELEASE_DIR};tar -czvf ${RELEASE_FILE}.tar.gz ${RELEASE_FILE}

create_release:
	github-release release \
    --user ${USER} \
    --repo ${REPO} \
    --tag ${VERSION} \
	--name "${VERSION}" \
    --description "[RELEASE NOTE](https://github.com/nginmesh/nginmesh/wiki/v${VERSION})"

upload_release:
	github-release upload \
    	--user ${USER} \
    	--repo ${REPO} \
    	--tag ${VERSION} \
    	--name "${RELEASE_FILE}.tar.gz" \
    	--file ${RELEASE_DIR}/${RELEASE_FILE}.tar.gz



delete_release:
	github-release delete \
    --user ${USER} \
    --repo ${REPO} \
    --tag ${VERSION}



clean:	clean_release
