dev: 
	air

tfdestroy: 
	terraform -chdir=./infra destroy 

tfapply: 
	terraform -chdir=./infra apply

