package symbol

type SymbolID int

const (
	SINGLE_INPUT       SymbolID = 256
	FILE_INPUT         SymbolID = 257
	EVAL_INPUT         SymbolID = 258
	DECORATOR          SymbolID = 259
	DECORATORS         SymbolID = 260
	DECORATED          SymbolID = 261
	ASYNC_FUNCDEF      SymbolID = 262
	FUNCDEF            SymbolID = 263
	PARAMETERS         SymbolID = 264
	TYPEDARGSLIST      SymbolID = 265
	TFPDEF             SymbolID = 266
	VARARGSLIST        SymbolID = 267
	VFPDEF             SymbolID = 268
	STMT               SymbolID = 269
	SIMPLE_STMT        SymbolID = 270
	SMALL_STMT         SymbolID = 271
	EXPR_STMT          SymbolID = 272
	TESTLIST_STAR_EXPR SymbolID = 273
	AUGASSIGN          SymbolID = 274
	DEL_STMT           SymbolID = 275
	PASS_STMT          SymbolID = 276
	FLOW_STMT          SymbolID = 277
	BREAK_STMT         SymbolID = 278
	CONTINUE_STMT      SymbolID = 279
	RETURN_STMT        SymbolID = 280
	YIELD_STMT         SymbolID = 281
	RAISE_STMT         SymbolID = 282
	IMPORT_STMT        SymbolID = 283
	IMPORT_NAME        SymbolID = 284
	IMPORT_FROM        SymbolID = 285
	IMPORT_AS_NAME     SymbolID = 286
	DOTTED_AS_NAME     SymbolID = 287
	IMPORT_AS_NAMES    SymbolID = 288
	DOTTED_AS_NAMES    SymbolID = 289
	DOTTED_NAME        SymbolID = 290
	GLOBAL_STMT        SymbolID = 291
	NONLOCAL_STMT      SymbolID = 292
	ASSERT_STMT        SymbolID = 293
	COMPOUND_STMT      SymbolID = 294
	ASYNC_STMT         SymbolID = 295
	IF_STMT            SymbolID = 296
	WHILE_STMT         SymbolID = 297
	FOR_STMT           SymbolID = 298
	TRY_STMT           SymbolID = 299
	WITH_STMT          SymbolID = 300
	WITH_ITEM          SymbolID = 301
	EXCEPT_CLAUSE      SymbolID = 302
	SUITE              SymbolID = 303
	TEST               SymbolID = 304
	TEST_NOCOND        SymbolID = 305
	LAMBDEF            SymbolID = 306
	LAMBDEF_NOCOND     SymbolID = 307
	OR_TEST            SymbolID = 308
	AND_TEST           SymbolID = 309
	NOT_TEST           SymbolID = 310
	COMPARISON         SymbolID = 311
	COMP_OP            SymbolID = 312
	STAR_EXPR          SymbolID = 313
	EXPR               SymbolID = 314
	XOR_EXPR           SymbolID = 315
	AND_EXPR           SymbolID = 316
	SHIFT_EXPR         SymbolID = 317
	ARITH_EXPR         SymbolID = 318
	TERM               SymbolID = 319
	FACTOR             SymbolID = 320
	POWER              SymbolID = 321
	ATOM_EXPR          SymbolID = 322
	ATOM               SymbolID = 323
	TESTLIST_COMP      SymbolID = 324
	TRAILER            SymbolID = 325
	SUBSCRIPTLIST      SymbolID = 326
	SUBSCRIPT          SymbolID = 327
	SLICEOP            SymbolID = 328
	EXPRLIST           SymbolID = 329
	TESTLIST           SymbolID = 330
	DICTORSETMAKER     SymbolID = 331
	CLASSDEF           SymbolID = 332
	ARGLIST            SymbolID = 333
	ARGUMENT           SymbolID = 334
	COMP_ITER          SymbolID = 335
	COMP_FOR           SymbolID = 336
	COMP_IF            SymbolID = 337
	ENCODING_DECL      SymbolID = 338
	YIELD_EXPR         SymbolID = 339
	YIELD_ARG          SymbolID = 340
)