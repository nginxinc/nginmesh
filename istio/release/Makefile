VERSION=0.2.12-RC2
REPO=nginmesh
USER=nginmesh
RELEASE_FILE=nginmesh.tar.gz


update-install:
	./updateVersion.sh


release:	release_tar create_release upload_release

clean_release:  
	rm -rf target

release_tar:    clean_release
	mkdir -p target
	cp -r samples target
	cp -r install target
	rm -rf target/install/kubernetes/templates
	rm -rf target/install/kubernetes/.gitignore
	rm -rf target/install/kubernetes/*.bak
	cd target;tar -czvf ${RELEASE_FILE} install samples            

create_release:
	github-release release \
    --user ${USER} \
    --repo ${REPO} \
    --tag ${VERSION} \
    --name "${VERSION}" \
    --pre-release

upload_release:
	github-release upload \
    	--user ${USER} \
    	--repo ${REPO} \
    	--tag ${VERSION} \
    	--name "${RELEASE_FILE}" \
    	--file target/${RELEASE_FILE}



delete_release:
	github-release delete \
    --user ${USER} \
    --repo ${REPO} \
    --tag ${VERSION}



clean:	clean_release
