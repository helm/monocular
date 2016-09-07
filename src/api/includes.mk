define check-static-binary
	  if file $(1) | egrep -q "(statically linked|Mach-O)"; then \
	    echo ""; \
	  else \
	    echo "The binary file $(1) is not statically linked. Build canceled"; \
	    exit 1; \
	  fi
endef
