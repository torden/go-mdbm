#ifndef H_MDBM_BINDING_ONCE

#define H_MDBM_BINDING_ONCE

#include <stdlib.h>
#include <stdio.h>
#include <sys/syslog.h>
#include <unistd.h>
#include <string.h>
#include <time.h>
#include <mdbm.h>
#include <mdbm_log.h>

#define LF_SKIP         -1

#define LT_NONE         0
#define LT_LOCK         1
#define LT_SMART        2
#define LT_SHARED       3
#define LT_PLOCK        5
#define LT_TRY_LOCK     11
#define LT_TRY_SMART    12
#define LT_TRY_SHARED   13
#define LT_TRY_PLOCK    15

extern void get_mdbm_iter(MDBM_ITER *iter);
extern mdbm_ubig_t get_pageno_of_mdbm_iter(MDBM_ITER *iter);
extern int get_next_of_mdbm_iter(MDBM_ITER *iter);
extern int set_mdbm_store_with_lock(MDBM *db, datum key, datum val, int flags, int locktype, int lockflags);
extern int set_mdbm_store_r_with_lock(MDBM *db, datum *key, datum *val, int flags, MDBM_ITER *iter,  int locktype, int lockflags);
extern int set_mdbm_store_str_with_lock(MDBM *db, const char *key, const char *val, int flags, int locktype, int lockflags);
extern char *get_mdbm_fetch_with_lock(MDBM *db, datum key, int locktype, int lockflags);
extern int get_mdbm_fetch_r_with_lock(MDBM *db, datum *key, datum *val, MDBM_ITER *iter, int locktype, int lockflags);
extern int get_mdbm_fetch_dup_r_with_lock(MDBM *db, datum *key, datum *val, MDBM_ITER *iter, int locktype, int lockflags);
extern char *get_mdbm_fetch_str_with_lock(MDBM *db, const char *key, int locktype, int lockflags);
extern int set_mdbm_delete_with_lock(MDBM *db, datum key, int locktype, int lockflags);
extern int set_mdbm_delete_r_with_lock(MDBM *db, datum key, MDBM_ITER *iter, int locktype, int lockflags);
extern int set_mdbm_delete_str_with_lock(MDBM *db, const char *key, int locktype, int lockflags);
extern int dummy_clean_func(MDBM *, const datum*, const datum*, struct mdbm_clean_data *, int* quit); 
#endif
