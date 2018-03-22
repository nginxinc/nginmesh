
cat << END                                                    >> $NGX_MAKEFILE

cargo:
	cargo build --release --manifest-path $ngx_addon_dir/../Cargo.toml --lib --all

END
