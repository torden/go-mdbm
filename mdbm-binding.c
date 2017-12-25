#include <mdbm-binding.h>

extern void get_mdbm_iter(MDBM_ITER *iter) {
	MDBM_ITER_INIT(iter);
}

extern mdbm_ubig_t get_pageno_of_mdbm_iter(MDBM_ITER *iter) {
	return (*iter).m_pageno;
}

extern int get_next_of_mdbm_iter(MDBM_ITER *iter) {
	return (*iter).m_next;
}

//debug
static inline void echo_str(const char *str) {
	fprintf(stderr,"----- echo_str -----\n");
	fprintf(stderr,"[%s][%d]\n", str, (int)strlen(str));
	fprintf(stderr,"--------------------\n");
}

static inline const char *get_locktype_name(int locktype) {

    switch(locktype) {
        case LT_SMART:
            return (const char *)"SMART";
        case LT_SHARED:
            return (const char *)"SHARED";
        case LT_PLOCK:
            return (const char *)"PLOCK";
        case LT_TRY_LOCK:
            return (const char *)"TRY_LOCK";
        case LT_TRY_SMART:
            return (const char *)"TRY_LOCK_SMART";
        case LT_TRY_SHARED:
            return (const char *)"TRY_LOCK_SHARED";
        case LT_TRY_PLOCK:
            return (const char *)"TRY_PLOCK";
        case LT_NONE:
            return (const char *)"NONE";
        default: //with case LT_LOCK:
            return (const char *)"LOCK";
            break;
    }
}


//LOCK
static inline int common_lock_func(MDBM *db, datum *key, int locktype, int lockflags) {

    int rv = -1;

    if(locktype == LT_NONE) {
        return 1;
    }

    if(key == NULL && (locktype == LT_SMART || locktype == LT_PLOCK || locktype == LT_TRY_SMART || locktype == LT_TRY_PLOCK)) {
        fprintf(stderr, "Not support Lock(=%s) without key\n", get_locktype_name(locktype));
        return rv;
    }

    switch(locktype) {
        case LT_SMART:
            rv = mdbm_lock_smart(db, key, lockflags);
            break;
        case LT_SHARED:
            rv = mdbm_lock_shared(db);
            break;
        case LT_PLOCK:
            rv = mdbm_plock(db, key, lockflags);
            break;
        case LT_TRY_LOCK:
            rv = mdbm_trylock(db);
            break;
        case LT_TRY_SMART:
            rv = mdbm_trylock_smart(db, key, lockflags);
            break;
        case LT_TRY_SHARED:
            rv = mdbm_trylock_shared(db);
            break;
        case LT_TRY_PLOCK:
            rv = mdbm_plock(db, key, lockflags);
            break;
        case LT_LOCK:
            rv = mdbm_lock(db);
            break;
        default:
            fprintf(stderr, "Not support Lock(=%s,%d) \n", get_locktype_name(locktype), locktype);
            break;

    }

    return rv;
}

//UNLOCK
static inline int common_unlock_func(MDBM *db, datum *key, int locktype, int lockflags) {

    int rv = -1;

    if(locktype == LT_NONE) {
        return 1;
    }

    if(key == NULL && (locktype == LT_SMART || locktype == LT_PLOCK || locktype == LT_TRY_SMART || locktype == LT_TRY_PLOCK)) {
        return rv;
    }


    switch(locktype) {
        case LT_SMART:
            rv = mdbm_unlock_smart(db, key, lockflags);
            break;
        case LT_PLOCK:
            rv = mdbm_punlock(db, key, lockflags);
            break;
        case LT_TRY_SMART:
            rv = mdbm_unlock_smart(db, key, lockflags);
            break;
        case LT_TRY_PLOCK:
            rv = mdbm_punlock(db, key, lockflags);
            break;
        case LT_LOCK:
            rv = mdbm_unlock(db);
            break;
        default:
            rv = mdbm_unlock(db);
            break;
    }

    return rv;
}

// STORE
extern int set_mdbm_store_with_lock(MDBM *db, datum key, datum val, int flags, int locktype, int lockflags) {

	int rv;

    rv = common_lock_func(db, &key, locktype, lockflags);
    if(rv != 1) {
        return rv;
    }

    rv = mdbm_store(db,key,val,flags);

    common_unlock_func(db, &key, locktype, lockflags);
    return rv;
}

// STORE_R
extern int set_mdbm_store_r_with_lock(MDBM *db, datum *key, datum *val, int flags, MDBM_ITER *iter,  int locktype, int lockflags) {

	int rv;
    rv = common_lock_func(db, key, locktype, lockflags);
    if(rv != 1) {
        return rv;
    }

	rv = mdbm_store_r(db,key,val,flags, iter);

    common_unlock_func(db, key, locktype, lockflags);
	return rv;
}


///STORE_STR
extern int set_mdbm_store_str_with_lock(MDBM *db, const char *key, const char *val, int flags, int locktype, int lockflags) {

	int rv;
	datum lockkey = {(char *)key, strlen(key)+1};

    rv = common_lock_func(db, &lockkey, locktype, lockflags);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_store_str(db, key, val, flags);

    common_unlock_func(db, &lockkey, locktype, lockflags);
	return rv;
}

// FETCH
extern char *get_mdbm_fetch_with_lock(MDBM *db, datum key, int locktype, int lockflags) {

	int rv;
	datum val;
	char *buf = NULL;

    rv = common_lock_func(db, &key, locktype, lockflags);
	if(rv != 1) {
		return buf;
	}

	val = mdbm_fetch(db, key);
	if (val.dptr != NULL) {
		buf = (char *)calloc(val.dsize+1, sizeof(char));
		if (buf == NULL) {
			perror("failed allocation:");
		} else {
			strncpy(buf, val.dptr, val.dsize);
		}
	}

    common_unlock_func(db, &key, locktype, lockflags);
	return buf;
}

// FETCH_R
extern int get_mdbm_fetch_r_with_lock(MDBM *db, datum *key, datum *val, MDBM_ITER *iter, int locktype, int lockflags) {

	int rv;
    rv = common_lock_func(db, key, locktype, lockflags);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_fetch_r(db, key, val, iter);

    common_unlock_func(db, key, locktype, lockflags);
	return rv;
}

// FETCH_DUP_R
extern int get_mdbm_fetch_dup_r_with_lock(MDBM *db, datum *key, datum *val, MDBM_ITER *iter, int locktype, int lockflags) {

	int rv;
    rv = common_lock_func(db, key, locktype, lockflags);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_fetch_dup_r(db, key, val, iter);

    common_unlock_func(db, key, locktype, lockflags);
	return rv;
}



// FETCH_STR
extern char *get_mdbm_fetch_str_with_lock(MDBM *db, const char *key, int locktype, int lockflags) {

	int rv;
	char *retval = NULL;

	datum lockkey = {(char *)key, strlen(key)+1};
    rv = common_lock_func(db, &lockkey, locktype, lockflags);
	if(rv != 1) {
		return NULL;
	}

	retval = mdbm_fetch_str(db, key);
    common_unlock_func(db, &lockkey, locktype, lockflags);
	return retval;
}

// DELETE
extern int set_mdbm_delete_with_lock(MDBM *db, datum key, int locktype, int lockflags) {

    int rv;
    rv = common_lock_func(db, &key, locktype, lockflags);
	if(rv != 1) {
		return rv;
	}

    rv = mdbm_delete(db, key);

    common_unlock_func(db, &key, locktype, lockflags);
	return rv;
}

// DELETE_R
extern int set_mdbm_delete_r_with_lock(MDBM *db, datum key, MDBM_ITER *iter, int locktype, int lockflags) {

    int rv;

    datum *lockkey = NULL;
    if(key.dsize >= 0) {
        lockkey = &key;
    }
       
    rv = common_lock_func(db, lockkey, locktype, lockflags);
	if(rv != 1) {
		return rv;
	}

    rv = mdbm_delete_r(db, iter);

    common_unlock_func(db, lockkey, locktype, lockflags);
	return rv;
}

extern int set_mdbm_delete_str_with_lock(MDBM *db, const char *key, int locktype, int lockflags) {

	int rv;

	datum lockkey = {(char *)key, strlen(key)};
    rv = common_lock_func(db, &lockkey, locktype, lockflags);
	if(rv != 1) {
		return rv;
	}

	rv = mdbm_delete_str(db, key);

    common_unlock_func(db, &lockkey, locktype, lockflags);
	return rv;
}

extern int dummy_clean_func(MDBM *, const datum*, const datum*, struct mdbm_clean_data *, int* quit) {
    *quit = 0;
    return 1;
}
