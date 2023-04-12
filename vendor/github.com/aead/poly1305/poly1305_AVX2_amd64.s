// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// This code is inspired by the poly1305 AVX2 implementation by Shay Gueron, and Vlad Krasnov.

// +build amd64, !gccgo, !appengine

#include "textflag.h"

DATA addMaskAVX2<>+0x00(SB)/8, $0x3FFFFFF
DATA addMaskAVX2<>+0x08(SB)/8, $0x3FFFFFF
DATA addMaskAVX2<>+0x10(SB)/8, $0x3FFFFFF
DATA addMaskAVX2<>+0x18(SB)/8, $0x3FFFFFF
GLOBL addMaskAVX2<>(SB), RODATA, $32

DATA poly1305MaskAVX2<>+0x00(SB)/8, $0xFFFFFFC0FFFFFFF
DATA poly1305MaskAVX2<>+0x08(SB)/8, $0xFFFFFFC0FFFFFFF
DATA poly1305MaskAVX2<>+0x10(SB)/8, $0xFFFFFFC0FFFFFFF
DATA poly1305MaskAVX2<>+0x18(SB)/8, $0xFFFFFFC0FFFFFFF
DATA poly1305MaskAVX2<>+0x20(SB)/8, $0xFFFFFFC0FFFFFFC
DATA poly1305MaskAVX2<>+0x28(SB)/8, $0xFFFFFFC0FFFFFFC
DATA poly1305MaskAVX2<>+0x30(SB)/8, $0xFFFFFFC0FFFFFFC
DATA poly1305MaskAVX2<>+0x38(SB)/8, $0xFFFFFFC0FFFFFFC
GLOBL poly1305MaskAVX2<>(SB), RODATA, $64

DATA oneBit<>+0x00(SB)/8, $0x1000000
DATA oneBit<>+0x08(SB)/8, $0x1000000
DATA oneBit<>+0x10(SB)/8, $0x1000000
DATA oneBit<>+0x18(SB)/8, $0x1000000
GLOBL oneBit<>(SB), RODATA, $32

DATA fixPermutation<>+0x00(SB)/4, $6
DATA fixPermutation<>+0x04(SB)/4, $7
DATA fixPermutation<>+0x08(SB)/4, $6
DATA fixPermutation<>+0x0c(SB)/4, $7
DATA fixPermutation<>+0x10(SB)/4, $6
DATA fixPermutation<>+0x14(SB)/4, $7
DATA fixPermutation<>+0x18(SB)/4, $6
DATA fixPermutation<>+0x1c(SB)/4, $7
DATA fixPermutation<>+0x20(SB)/4, $4
DATA fixPermutation<>+0x24(SB)/4, $5
DATA fixPermutation<>+0x28(SB)/4, $6
DATA fixPermutation<>+0x2c(SB)/4, $7
DATA fixPermutation<>+0x30(SB)/4, $6
DATA fixPermutation<>+0x34(SB)/4, $7
DATA fixPermutation<>+0x38(SB)/4, $6
DATA fixPermutation<>+0x3c(SB)/4, $7
DATA fixPermutation<>+0x40(SB)/4, $2
DATA fixPermutation<>+0x44(SB)/4, $3
DATA fixPermutation<>+0x48(SB)/4, $6
DATA fixPermutation<>+0x4c(SB)/4, $7
DATA fixPermutation<>+0x50(SB)/4, $4
DATA fixPermutation<>+0x54(SB)/4, $5
DATA fixPermutation<>+0x58(SB)/4, $6
DATA fixPermutation<>+0x5c(SB)/4, $7
DATA fixPermutation<>+0x60(SB)/4, $0
DATA fixPermutation<>+0x64(SB)/4, $1
DATA fixPermutation<>+0x68(SB)/4, $4
DATA fixPermutation<>+0x6c(SB)/4, $5
DATA fixPermutation<>+0x70(SB)/4, $2
DATA fixPermutation<>+0x74(SB)/4, $3
DATA fixPermutation<>+0x78(SB)/4, $6
DATA fixPermutation<>+0x7c(SB)/4, $7
GLOBL fixPermutation<>(SB), RODATA, $128

TEXT ·initializeAVX2(SB), $0-16
	MOVQ state+0(FP), DI
	MOVQ key+8(FP), SI

	MOVQ $addMaskAVX2<>(SB), R8

	MOVOU 16*1(SI), X10
	MOVOU X10, 288(DI)
	PXOR  X10, X10
	MOVOU X10, 304(DI)

	MOVD X10, 320(DI)
	MOVQ 8*0(SI), X5
	MOVQ 8*1(SI), X10

	VZEROUPPER

	MOVQ  $poly1305MaskAVX2<>(SB), R9
	VPAND (R9), X5, X5
	VPAND 32(R9), X10, X10

	VMOVDQU 0(R8), X0
	VPSRLQ  $26, X5, X6
	VPAND   X0, X5, X5
	VPSRLQ  $26, X6, X7
	VPAND   X0, X6, X6
	VPSLLQ  $12, X10, X11
	VPXOR   X11, X7, X7
	VPSRLQ  $26, X7, X8
	VPSRLQ  $40, X10, X9
	VPAND   X0, X7, X7
	VPAND   X0, X8, X8

	BYTE $0xc5; BYTE $0xd1; BYTE $0xf4; BYTE $0xc5             // VPMULUDQ	X5, X5, X0
	BYTE $0xc5; BYTE $0xd1; BYTE $0xf4; BYTE $0xce             // VPMULUDQ	X6, X5, X1
	BYTE $0xc5; BYTE $0xd1; BYTE $0xf4; BYTE $0xd7             // VPMULUDQ	X7, X5, X2
	BYTE $0xc4; BYTE $0xc1; BYTE $0x51; BYTE $0xf4; BYTE $0xd8 // VPMULUDQ	X8, X5, X3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x51; BYTE $0xf4; BYTE $0xe1 // VPMULUDQ	X9, X5, X4

	VPSLLQ $1, X1, X1
	VPSLLQ $1, X2, X2
	BYTE   $0xc5; BYTE $0x49; BYTE $0xf4; BYTE $0xd6             // VPMULUDQ X6, X6, X10
	VPADDQ X10, X2, X2
	BYTE   $0xc5; BYTE $0x49; BYTE $0xf4; BYTE $0xd7             // VPMULUDQ X7, X6, X10
	VPADDQ X10, X3, X3
	BYTE   $0xc4; BYTE $0x41; BYTE $0x49; BYTE $0xf4; BYTE $0xd0 // VPMULUDQ X8, X6, X10
	VPADDQ X10, X4, X4
	BYTE   $0xc4; BYTE $0x41; BYTE $0x49; BYTE $0xf4; BYTE $0xe1 // VPMULUDQ X9, X6, X12
	VPSLLQ $1, X3, X3
	VPSLLQ $1, X4, X4
	BYTE   $0xc5; BYTE $0x41; BYTE $0xf4; BYTE $0xd7             // VPMULUDQ X7, X7, X10
	VPADDQ X10, X4, X4
	BYTE   $0xc4; BYTE $0x41; BYTE $0x41; BYTE $0xf4; BYTE $0xd0 // VPMULUDQ X8, X7, X10
	VPADDQ X10, X12, X12
	BYTE   $0xc4; BYTE $0x41; BYTE $0x41; BYTE $0xf4; BYTE $0xe9 // VPMULUDQ X9, X7, X13
	VPSLLQ $1, X12, X12
	VPSLLQ $1, X13, X13
	BYTE   $0xc4; BYTE $0x41; BYTE $0x39; BYTE $0xf4; BYTE $0xd0 // VPMULUDQ X8, X8, X10
	VPADDQ X10, X13, X13
	BYTE   $0xc4; BYTE $0x41; BYTE $0x39; BYTE $0xf4; BYTE $0xf1 // VPMULUDQ X9, X8, X14
	VPSLLQ $1, X14, X14
	BYTE   $0xc4; BYTE $0x41; BYTE $0x31; BYTE $0xf4; BYTE $0xf9 // VPMULUDQ X9, X9, X15

	VPSRLQ $26, X4, X10
	VPAND  0(R8), X4, X4
	VPADDQ X10, X12, X12

	VPSLLQ $2, X12, X10
	VPADDQ X10, X12, X12
	VPSLLQ $2, X13, X10
	VPADDQ X10, X13, X13
	VPSLLQ $2, X14, X10
	VPADDQ X10, X14, X14
	VPSLLQ $2, X15, X10
	VPADDQ X10, X15, X15

	VPADDQ X12, X0, X0
	VPADDQ X13, X1, X1
	VPADDQ X14, X2, X2
	VPADDQ X15, X3, X3

	VPSRLQ $26, X0, X10
	VPAND  0(R8), X0, X0
	VPADDQ X10, X1, X1
	VPSRLQ $26, X1, X10
	VPAND  0(R8), X1, X1
	VPADDQ X10, X2, X2
	VPSRLQ $26, X2, X10
	VPAND  0(R8), X2, X2
	VPADDQ X10, X3, X3
	VPSRLQ $26, X3, X10
	VPAND  0(R8), X3, X3
	VPADDQ X10, X4, X4

	BYTE $0xc5; BYTE $0xf9; BYTE $0x6c; BYTE $0xed             // VPUNPCKLQDQ	X5, X0, X5
	BYTE $0xc5; BYTE $0xf1; BYTE $0x6c; BYTE $0xf6             // VPUNPCKLQDQ	X6, X1, X6
	BYTE $0xc5; BYTE $0xe9; BYTE $0x6c; BYTE $0xff             // VPUNPCKLQDQ	X7, X2, X7
	BYTE $0xc4; BYTE $0x41; BYTE $0x61; BYTE $0x6c; BYTE $0xc0 // VPUNPCKLQDQ	X8, X3, X8
	BYTE $0xc4; BYTE $0x41; BYTE $0x59; BYTE $0x6c; BYTE $0xc9 // VPUNPCKLQDQ	X9, X4, X9

	VMOVDQU X5, 0+16(DI)
	VMOVDQU X6, 32+16(DI)
	VMOVDQU X7, 64+16(DI)
	VMOVDQU X8, 96+16(DI)
	VMOVDQU X9, 128+16(DI)

	VPSLLQ $2, X6, X1
	VPSLLQ $2, X7, X2
	VPSLLQ $2, X8, X3
	VPSLLQ $2, X9, X4

	VPADDQ X1, X6, X1
	VPADDQ X2, X7, X2
	VPADDQ X3, X8, X3
	VPADDQ X4, X9, X4

	VMOVDQU X1, 160+16(DI)
	VMOVDQU X2, 192+16(DI)
	VMOVDQU X3, 224+16(DI)
	VMOVDQU X4, 256+16(DI)

	VPSHUFD $68, X5, X0
	VPSHUFD $68, X6, X1
	VPSHUFD $68, X7, X2
	VPSHUFD $68, X8, X3
	VPSHUFD $68, X9, X4

	VMOVDQU 0+16(DI), X10
	BYTE    $0xc4; BYTE $0xc1; BYTE $0x79; BYTE $0xf4; BYTE $0xea // VPMULUDQ	X10, X0, X5
	BYTE    $0xc4; BYTE $0xc1; BYTE $0x71; BYTE $0xf4; BYTE $0xf2 // VPMULUDQ	X10, X1, X6
	BYTE    $0xc4; BYTE $0xc1; BYTE $0x69; BYTE $0xf4; BYTE $0xfa // VPMULUDQ	X10, X2, X7
	BYTE    $0xc4; BYTE $0x41; BYTE $0x61; BYTE $0xf4; BYTE $0xc2 // VPMULUDQ	X10, X3, X8
	BYTE    $0xc4; BYTE $0x41; BYTE $0x59; BYTE $0xf4; BYTE $0xca // VPMULUDQ	X10, X4, X9

	VMOVDQU 160+16(DI), X10
	BYTE    $0xc4; BYTE $0x41; BYTE $0x59; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10 ,X4, X11
	VPADDQ  X11, X5, X5

	VMOVDQU 32+16(DI), X10
	BYTE    $0xc4; BYTE $0x41; BYTE $0x79; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10 ,X0, X11
	VPADDQ  X11, X6, X6
	BYTE    $0xc4; BYTE $0x41; BYTE $0x71; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10 ,X1, X11
	VPADDQ  X11, X7, X7
	BYTE    $0xc4; BYTE $0x41; BYTE $0x69; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10 ,X2, X11
	VPADDQ  X11, X8, X8
	BYTE    $0xc4; BYTE $0x41; BYTE $0x61; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10 ,X3, X11
	VPADDQ  X11, X9, X9

	VMOVDQU 192+16(DI), X10
	BYTE    $0xc4; BYTE $0x41; BYTE $0x61; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10 ,X3, X11
	VPADDQ  X11, X5, X5
	BYTE    $0xc4; BYTE $0x41; BYTE $0x59; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10 ,X4, X11
	VPADDQ  X11, X6, X6

	VMOVDQU 64+16(DI), X10
	BYTE    $0xc4; BYTE $0x41; BYTE $0x79; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10 ,X0, X11
	VPADDQ  X11, X7, X7
	BYTE    $0xc4; BYTE $0x41; BYTE $0x71; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10 ,X1, X11
	VPADDQ  X11, X8, X8
	BYTE    $0xc4; BYTE $0x41; BYTE $0x69; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10 ,X2, X11
	VPADDQ  X11, X9, X9

	VMOVDQU 224+16(DI), X10

	BYTE   $0xc4; BYTE $0x41; BYTE $0x69; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10, X2, X11
	VPADDQ X11, X5, X5
	BYTE   $0xc4; BYTE $0x41; BYTE $0x61; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10, X3, X11
	VPADDQ X11, X6, X6
	BYTE   $0xc4; BYTE $0x41; BYTE $0x59; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10, X4, X11
	VPADDQ X11, X7, X7

	VMOVDQU 96+16(DI), X10
	BYTE    $0xc4; BYTE $0x41; BYTE $0x79; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10, X0, X11
	VPADDQ  X11, X8, X8
	BYTE    $0xc4; BYTE $0x41; BYTE $0x71; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10, X1, X11
	VPADDQ  X11, X9, X9

	VMOVDQU 256+16(DI), X10
	BYTE    $0xc4; BYTE $0x41; BYTE $0x71; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10, X1, X11
	VPADDQ  X11, X5, X5
	BYTE    $0xc4; BYTE $0x41; BYTE $0x69; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10, X2, X11
	VPADDQ  X11, X6, X6
	BYTE    $0xc4; BYTE $0x41; BYTE $0x61; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10, X3, X11
	VPADDQ  X11, X7, X7
	BYTE    $0xc4; BYTE $0x41; BYTE $0x59; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10, X4, X11
	VPADDQ  X11, X8, X8

	VMOVDQU 128+16(DI), X10
	BYTE    $0xc4; BYTE $0x41; BYTE $0x79; BYTE $0xf4; BYTE $0xda // VPMULUDQ	X10, X0, X11
	VPADDQ  X11, X9, X9

	VMOVDQU 0(R8), X12

	VPSRLQ $26, X8, X10
	VPADDQ X10, X9, X9
	VPAND  X12, X8, X8
	VPSRLQ $26, X9, X10
	VPSLLQ $2, X10, X11
	VPADDQ X11, X10, X10
	VPADDQ X10, X5, X5
	VPAND  X12, X9, X9
	VPSRLQ $26, X5, X10
	VPAND  X12, X5, X5
	VPADDQ X10, X6, X6
	VPSRLQ $26, X6, X10
	VPAND  X12, X6, X6
	VPADDQ X10, X7, X7
	VPSRLQ $26, X7, X10
	VPAND  X12, X7, X7
	VPADDQ X10, X8, X8
	VPSRLQ $26, X8, X10
	VPAND  X12, X8, X8
	VPADDQ X10, X9, X9

	VMOVDQU X5, 0(DI)
	VMOVDQU X6, 32(DI)
	VMOVDQU X7, 64(DI)
	VMOVDQU X8, 96(DI)
	VMOVDQU X9, 128(DI)

	VPSLLQ $2, X6, X1
	VPSLLQ $2, X7, X2
	VPSLLQ $2, X8, X3
	VPSLLQ $2, X9, X4

	VPADDQ X1, X6, X1
	VPADDQ X2, X7, X2
	VPADDQ X3, X8, X3
	VPADDQ X4, X9, X4

	VMOVDQU X1, 160(DI)
	VMOVDQU X2, 192(DI)
	VMOVDQU X3, 224(DI)
	VMOVDQU X4, 256(DI)

	RET

TEXT ·updateAVX2(SB), $0-24
	MOVQ state+0(FP), DI
	MOVQ msg+8(FP), SI
	MOVQ msg_len+16(FP), DX

	MOVD 304(DI), X0
	MOVD 308(DI), X1
	MOVD 312(DI), X2
	MOVD 316(DI), X3
	MOVD 320(DI), X4

	MOVQ $addMaskAVX2<>(SB), R12
	MOVQ $oneBit<>(SB), R13
	MOVQ $fixPermutation<>(SB), R15
	VZEROUPPER

	VMOVDQA (R12), Y12

	CMPQ DX, $128
	JB   BETWEEN_0_AND_128

AT_LEAST_128:
	VMOVDQU 32*0(SI), Y9
	VMOVDQU 32*1(SI), Y10
	ADDQ    $64, SI

	BYTE $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x6c; BYTE $0xfa             // VPUNPCKLQDQ	Y10,Y9,Y7
	BYTE $0xc4; BYTE $0x41; BYTE $0x35; BYTE $0x6d; BYTE $0xc2             // VPUNPCKHQDQ	Y10,Y9,Y8
	BYTE $0xc4; BYTE $0xe3; BYTE $0xfd; BYTE $0x00; BYTE $0xff; BYTE $0xd8 // VPERMQ	$216,Y7,Y7
	BYTE $0xc4; BYTE $0x43; BYTE $0xfd; BYTE $0x00; BYTE $0xc0; BYTE $0xd8 // VPERMQ	$216,Y8,Y8

	VPSRLQ $26, Y7, Y9
	VPAND  Y12, Y7, Y7
	VPADDQ Y7, Y0, Y0

	VPSRLQ $26, Y9, Y7
	VPAND  Y12, Y9, Y9
	VPADDQ Y9, Y1, Y1

	VPSLLQ $12, Y8, Y9
	VPXOR  Y9, Y7, Y7
	VPAND  Y12, Y7, Y7
	VPADDQ Y7, Y2, Y2

	VPSRLQ $26, Y9, Y7
	VPSRLQ $40, Y8, Y9
	VPAND  Y12, Y7, Y7
	VPXOR  (R13), Y9, Y9
	VPADDQ Y7, Y3, Y3
	VPADDQ Y9, Y4, Y4

	BYTE $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x59; BYTE $0x2f // VPBROADCASTQ	0(DI),  Y5
	BYTE $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xfd             // VPMULUDQ		Y5, Y0, Y7
	BYTE $0xc5; BYTE $0x75; BYTE $0xf4; BYTE $0xc5             // VPMULUDQ		Y5, Y1, Y8
	BYTE $0xc5; BYTE $0x6d; BYTE $0xf4; BYTE $0xcd             // VPMULUDQ		Y5, Y2, Y9
	BYTE $0xc5; BYTE $0x65; BYTE $0xf4; BYTE $0xd5             // VPMULUDQ		Y5, Y3, Y10
	BYTE $0xc5; BYTE $0x5d; BYTE $0xf4; BYTE $0xdd             // VPMULUDQ		Y5, Y4, Y11

	BYTE   $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x59; BYTE $0xaf; BYTE $0xa0; BYTE $0x00; BYTE $0x00; BYTE $0x00 // VPBROADCASTQ	160(DI), Y5
	BYTE   $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y4, Y6
	VPADDQ Y6, Y7, Y7

	BYTE   $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x59; BYTE $0x6f; BYTE $0x20 // VPBROADCASTQ	32(DI), Y5
	BYTE   $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5                         // VPMULUDQ		Y5, Y0, Y6
	VPADDQ Y6, Y8, Y8
	BYTE   $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5                         // VPMULUDQ	Y5, Y1, Y6
	VPADDQ Y6, Y9, Y9
	BYTE   $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5                         // VPMULUDQ	Y5, Y2, Y6
	VPADDQ Y6, Y10, Y10
	BYTE   $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5                         // VPMULUDQ	Y5, Y3, Y6
	VPADDQ Y6, Y11, Y11

	BYTE   $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x59; BYTE $0xaf; BYTE $0xc0; BYTE $0x00; BYTE $0x00; BYTE $0x00 // VPBROADCASTQ	192(DI), Y5
	BYTE   $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y3, Y6
	VPADDQ Y6, Y7, Y7
	BYTE   $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y4, Y6
	VPADDQ Y6, Y8, Y8

	BYTE   $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x59; BYTE $0x6f; BYTE $0x40 // VPBROADCASTQ	64(DI), Y5
	BYTE   $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5                         // VPMULUDQ		Y5, Y0, Y6
	VPADDQ Y6, Y9, Y9
	BYTE   $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5                         // VPMULUDQ		Y5, Y1, Y6
	VPADDQ Y6, Y10, Y10
	BYTE   $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5                         // VPMULUDQ		Y5, Y2, Y6
	VPADDQ Y6, Y11, Y11

	BYTE   $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x59; BYTE $0xaf; BYTE $0xe0; BYTE $0x00; BYTE $0x00; BYTE $0x00 // VPBROADCASTQ	224(DI), Y5
	BYTE   $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y2, Y6
	VPADDQ Y6, Y7, Y7
	BYTE   $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y3, Y6
	VPADDQ Y6, Y8, Y8
	BYTE   $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y4, Y6
	VPADDQ Y6, Y9, Y9

	BYTE   $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x59; BYTE $0x6f; BYTE $0x60 // VPBROADCASTQ	96(DI), Y5
	BYTE   $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5                         // VPMULUDQ		Y5, Y0, Y6
	VPADDQ Y6, Y10, Y10
	BYTE   $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5                         // VPMULUDQ		Y5, Y1, Y6
	VPADDQ Y6, Y11, Y11

	BYTE   $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x59; BYTE $0xaf; BYTE $0x00; BYTE $0x01; BYTE $0x00; BYTE $0x00 // VPBROADCASTQ	256(DI), Y5
	BYTE   $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y1, Y6
	VPADDQ Y6, Y7, Y7
	BYTE   $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y2,Y6
	VPADDQ Y6, Y8, Y8
	BYTE   $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y3, Y6
	VPADDQ Y6, Y9, Y9
	BYTE   $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y4, Y6
	VPADDQ Y6, Y10, Y10

	BYTE   $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x59; BYTE $0xaf; BYTE $0x80; BYTE $0x00; BYTE $0x00; BYTE $0x00 // VPBROADCASTQ	128(DI),Y5
	BYTE   $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5                                                             // VPMULUDQ		Y5, Y0, Y6
	VPADDQ Y6, Y11, Y11

	VPSRLQ $26, Y10, Y5
	VPADDQ Y5, Y11, Y11
	VPAND  Y12, Y10, Y10

	VPSRLQ $26, Y11, Y5
	VPSLLQ $2, Y5, Y6
	VPADDQ Y6, Y5, Y5
	VPADDQ Y5, Y7, Y7
	VPAND  Y12, Y11, Y11

	VPSRLQ $26, Y7, Y5
	VPAND  Y12, Y7, Y0
	VPADDQ Y5, Y8, Y8
	VPSRLQ $26, Y8, Y5
	VPAND  Y12, Y8, Y1
	VPADDQ Y5, Y9, Y9
	VPSRLQ $26, Y9, Y5
	VPAND  Y12, Y9, Y2
	VPADDQ Y5, Y10, Y10
	VPSRLQ $26, Y10, Y5
	VPAND  Y12, Y10, Y3
	VPADDQ Y5, Y11, Y4

	SUBQ $64, DX
	CMPQ DX, $128
	JAE  AT_LEAST_128

BETWEEN_0_AND_128:
	CMPQ DX, $64
	JB   BETWEEN_0_AND_64

	VMOVDQU 32*0(SI), Y9
	VMOVDQU 32*1(SI), Y10
	ADDQ    $64, SI

	BYTE $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x6c; BYTE $0xfa             // VPUNPCKLQDQ	Y10, Y9, Y7
	BYTE $0xc4; BYTE $0x41; BYTE $0x35; BYTE $0x6d; BYTE $0xc2             // VPUNPCKHQDQ	Y10, Y9, Y8
	BYTE $0xc4; BYTE $0xe3; BYTE $0xfd; BYTE $0x00; BYTE $0xff; BYTE $0xd8 // VPERMQ		$216, Y7, Y7
	BYTE $0xc4; BYTE $0x43; BYTE $0xfd; BYTE $0x00; BYTE $0xc0; BYTE $0xd8 // VPERMQ		$216, Y8, Y8

	VPSRLQ $26, Y7, Y9
	VPAND  Y12, Y7, Y7
	VPADDQ Y7, Y0, Y0

	VPSRLQ $26, Y9, Y7
	VPAND  Y12, Y9, Y9
	VPADDQ Y9, Y1, Y1

	VPSLLQ $12, Y8, Y9
	VPXOR  Y9, Y7, Y7
	VPAND  Y12, Y7, Y7
	VPADDQ Y7, Y2, Y2

	VPSRLQ $26, Y9, Y7
	VPSRLQ $40, Y8, Y9
	VPAND  Y12, Y7, Y7
	VPXOR  (R13), Y9, Y9
	VPADDQ Y7, Y3, Y3
	VPADDQ Y9, Y4, Y4

	VMOVDQU 0(DI), Y5

	BYTE $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xfd // VPMULUDQ	Y5, Y0, Y7
	BYTE $0xc5; BYTE $0x75; BYTE $0xf4; BYTE $0xc5 // VPMULUDQ	Y5, Y1, Y8
	BYTE $0xc5; BYTE $0x6d; BYTE $0xf4; BYTE $0xcd // VPMULUDQ	Y5, Y2, Y9
	BYTE $0xc5; BYTE $0x65; BYTE $0xf4; BYTE $0xd5 // VPMULUDQ	Y5, Y3, Y10
	BYTE $0xc5; BYTE $0x5d; BYTE $0xf4; BYTE $0xdd // VPMULUDQ	Y5, Y4, Y11

	VMOVDQU 160(DI), Y5
	BYTE    $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y4, Y6
	VPADDQ  Y6, Y7, Y7
	VMOVDQU 32(DI), Y5
	BYTE    $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y0, Y6
	VPADDQ  Y6, Y8, Y8
	BYTE    $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y1, Y6
	VPADDQ  Y6, Y9, Y9
	BYTE    $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y2, Y6
	VPADDQ  Y6, Y10, Y10
	BYTE    $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y3, Y6
	VPADDQ  Y6, Y11, Y11

	VMOVDQU 192(DI), Y5
	BYTE    $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y3, Y6
	VPADDQ  Y6, Y7, Y7
	BYTE    $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y4, Y6
	VPADDQ  Y6, Y8, Y8

	VMOVDQU 64(DI), Y5
	BYTE    $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y0, Y6
	VPADDQ  Y6, Y9, Y9
	BYTE    $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y1, Y6
	VPADDQ  Y6, Y10, Y10
	BYTE    $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y2, Y6
	VPADDQ  Y6, Y11, Y11

	VMOVDQU 224(DI), Y5
	BYTE    $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y2, Y6
	VPADDQ  Y6, Y7, Y7
	BYTE    $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y3, Y6
	VPADDQ  Y6, Y8, Y8
	BYTE    $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y4, Y6
	VPADDQ  Y6, Y9, Y9

	VMOVDQU 96(DI), Y5
	BYTE    $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y0, Y6
	VPADDQ  Y6, Y10, Y10
	BYTE    $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y1, Y6
	VPADDQ  Y6, Y11, Y11

	VMOVDQU 256(DI), Y5
	BYTE    $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y1, Y6
	VPADDQ  Y6, Y7, Y7
	BYTE    $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y2, Y6
	VPADDQ  Y6, Y8, Y8
	BYTE    $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y3, Y6
	VPADDQ  Y6, Y9, Y9
	BYTE    $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y4, Y6
	VPADDQ  Y6, Y10, Y10

	VMOVDQU 128(DI), Y5
	BYTE    $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5 // VPMULUDQ	Y5, Y0, Y6
	VPADDQ  Y6, Y11, Y11

	VPSRLQ $26, Y10, Y5
	VPADDQ Y5, Y11, Y11
	VPAND  Y12, Y10, Y10
	VPSRLQ $26, Y11, Y5
	VPSLLQ $2, Y5, Y6
	VPADDQ Y6, Y5, Y5
	VPADDQ Y5, Y7, Y7
	VPAND  Y12, Y11, Y11
	VPSRLQ $26, Y7, Y5
	VPAND  Y12, Y7, Y0
	VPADDQ Y5, Y8, Y8
	VPSRLQ $26, Y8, Y5
	VPAND  Y12, Y8, Y1
	VPADDQ Y5, Y9, Y9
	VPSRLQ $26, Y9, Y5
	VPAND  Y12, Y9, Y2
	VPADDQ Y5, Y10, Y10
	VPSRLQ $26, Y10, Y5
	VPAND  Y12, Y10, Y3
	VPADDQ Y5, Y11, Y4

	VPSRLDQ $8, Y0, Y7
	VPSRLDQ $8, Y1, Y8
	VPSRLDQ $8, Y2, Y9
	VPSRLDQ $8, Y3, Y10
	VPSRLDQ $8, Y4, Y11

	VPADDQ Y7, Y0, Y0
	VPADDQ Y8, Y1, Y1
	VPADDQ Y9, Y2, Y2
	VPADDQ Y10, Y3, Y3
	VPADDQ Y11, Y4, Y4

	BYTE $0xc4; BYTE $0xe3; BYTE $0xfd; BYTE $0x00; BYTE $0xf8; BYTE $0xaa // VPERMQ	$170, Y0, Y7
	BYTE $0xc4; BYTE $0x63; BYTE $0xfd; BYTE $0x00; BYTE $0xc1; BYTE $0xaa // VPERMQ	$170, Y1, Y8
	BYTE $0xc4; BYTE $0x63; BYTE $0xfd; BYTE $0x00; BYTE $0xca; BYTE $0xaa // VPERMQ	$170, Y2, Y9
	BYTE $0xc4; BYTE $0x63; BYTE $0xfd; BYTE $0x00; BYTE $0xd3; BYTE $0xaa // VPERMQ	$170, Y3, Y10
	BYTE $0xc4; BYTE $0x63; BYTE $0xfd; BYTE $0x00; BYTE $0xdc; BYTE $0xaa // VPERMQ	$170, Y4, Y11

	VPADDQ Y7, Y0, Y0
	VPADDQ Y8, Y1, Y1
	VPADDQ Y9, Y2, Y2
	VPADDQ Y10, Y3, Y3
	VPADDQ Y11, Y4, Y4
	SUBQ   $64, DX

BETWEEN_0_AND_64:
	TESTQ DX, DX
	JZ    DONE

	BYTE $0xc5; BYTE $0xfa; BYTE $0x7e; BYTE $0xc0 // VMOVQ	X0, X0
	BYTE $0xc5; BYTE $0xfa; BYTE $0x7e; BYTE $0xc9 // VMOVQ	X1, X1
	BYTE $0xc5; BYTE $0xfa; BYTE $0x7e; BYTE $0xd2 // VMOVQ	X2, X2
	BYTE $0xc5; BYTE $0xfa; BYTE $0x7e; BYTE $0xdb // VMOVQ	X3, X3
	BYTE $0xc5; BYTE $0xfa; BYTE $0x7e; BYTE $0xe4 // VMOVQ	X4, X4

	MOVQ  (R13), BX
	MOVQ  SP, AX
	TESTQ $15, DX
	JZ    FULL_BLOCKS

	SUBQ    $64, SP
	VPXOR   Y7, Y7, Y7
	VMOVDQU Y7, (SP)
	VMOVDQU Y7, 32(SP)

	XORQ BX, BX

FLUSH_BUFFER:
	MOVB (SI)(BX*1), CX
	MOVB CX, (SP)(BX*1)
	INCQ BX
	CMPQ DX, BX
	JNE  FLUSH_BUFFER

	MOVB $1, (SP)(BX*1)
	XORQ BX, BX
	MOVQ SP, SI

FULL_BLOCKS:
	CMPQ DX, $16
	JA   AT_LEAST_16

	BYTE    $0xc5; BYTE $0xfa; BYTE $0x7e; BYTE $0x3e             // VMOVQ	8*0(SI), X7
	BYTE    $0xc5; BYTE $0x7a; BYTE $0x7e; BYTE $0x46; BYTE $0x08 // VMOVQ	8*1(SI), X8
	BYTE    $0xc4; BYTE $0x61; BYTE $0xf9; BYTE $0x6e; BYTE $0xf3 // VMOVQ	BX ,X14
	VMOVDQA (R15), Y13
	JMP     MULTIPLY

AT_LEAST_16:
	CMPQ    DX, $32
	JA      AT_LEAST_32
	VMOVDQU 16*0(SI), X9
	VMOVDQU 16*1(SI), X10

	BYTE    $0xc4; BYTE $0x41; BYTE $0x7a; BYTE $0x7e; BYTE $0x75; BYTE $0x00 // VMOVQ		(R13), X14
	BYTE    $0xc4; BYTE $0x63; BYTE $0x89; BYTE $0x22; BYTE $0xf3; BYTE $0x01 // VPINSRQ	$1,BX, X14, X14
	VMOVDQA 32(R15), Y13
	BYTE    $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x6c; BYTE $0xfa             // VPUNPCKLQDQ	Y10, Y9, Y7
	BYTE    $0xc4; BYTE $0x41; BYTE $0x35; BYTE $0x6d; BYTE $0xc2             // VPUNPCKHQDQ	Y10, Y9, Y8
	JMP     MULTIPLY

AT_LEAST_32:
	CMPQ    DX, $48
	JA      AT_LEAST_48
	VMOVDQU 32*0(SI), Y9
	VMOVDQU 32*1(SI), X10

	BYTE    $0xc4; BYTE $0x41; BYTE $0x7a; BYTE $0x7e; BYTE $0x75; BYTE $0x00 // VMOVQ		0(R13), X14
	BYTE    $0xc4; BYTE $0x63; BYTE $0x89; BYTE $0x22; BYTE $0xf3; BYTE $0x01 // VPINSRQ	$1, BX, X14, X14
	BYTE    $0xc4; BYTE $0x43; BYTE $0xfd; BYTE $0x00; BYTE $0xf6; BYTE $0xc4 // VPERMQ		$196, Y14, Y14
	VMOVDQA 64(R15), Y13
	BYTE    $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x6c; BYTE $0xfa             // VPUNPCKLQDQ	Y10, Y9, Y7
	BYTE    $0xc4; BYTE $0x41; BYTE $0x35; BYTE $0x6d; BYTE $0xc2             // VPUNPCKHQDQ	Y10, Y9, Y8
	JMP     MULTIPLY

AT_LEAST_48:
	VMOVDQU 32*0(SI), Y9
	VMOVDQU 32*1(SI), Y10

	BYTE    $0xc4; BYTE $0x41; BYTE $0x7a; BYTE $0x7e; BYTE $0x75; BYTE $0x00 // VMOVQ	(R13),X14
	BYTE    $0xc4; BYTE $0x63; BYTE $0x89; BYTE $0x22; BYTE $0xf3; BYTE $0x01 // VPINSRQ	$1,BX,X14,X14
	BYTE    $0xc4; BYTE $0x43; BYTE $0xfd; BYTE $0x00; BYTE $0xf6; BYTE $0x40 // VPERMQ	$64,Y14,Y14
	VMOVDQA 96(R15), Y13
	BYTE    $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x6c; BYTE $0xfa             // VPUNPCKLQDQ	Y10, Y9, Y7
	BYTE    $0xc4; BYTE $0x41; BYTE $0x35; BYTE $0x6d; BYTE $0xc2             // VPUNPCKHQDQ	Y10, Y9, Y8

MULTIPLY:
	MOVQ AX, SP

	VPSRLQ $26, Y7, Y9
	VPAND  Y12, Y7, Y7
	VPADDQ Y7, Y0, Y0

	VPSRLQ $26, Y9, Y7
	VPAND  Y12, Y9, Y9
	VPADDQ Y9, Y1, Y1

	VPSLLQ $12, Y8, Y9
	VPXOR  Y9, Y7, Y7
	VPAND  Y12, Y7, Y7
	VPADDQ Y7, Y2, Y2

	VPSRLQ $26, Y9, Y7
	VPSRLQ $40, Y8, Y9
	VPAND  Y12, Y7, Y7
	VPXOR  Y14, Y9, Y9
	VPADDQ Y7, Y3, Y3
	VPADDQ Y9, Y4, Y4

	VMOVDQU 0(DI), Y5

	BYTE $0xc4; BYTE $0xe2; BYTE $0x15; BYTE $0x36; BYTE $0xed // VPERMD	Y5, Y13, Y5
	BYTE $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xfd             // VPMULUDQ	Y5, Y0, Y7
	BYTE $0xc5; BYTE $0x75; BYTE $0xf4; BYTE $0xc5             // VPMULUDQ	Y5, Y1, Y8
	BYTE $0xc5; BYTE $0x6d; BYTE $0xf4; BYTE $0xcd             // VPMULUDQ	Y5, Y2, Y9
	BYTE $0xc5; BYTE $0x65; BYTE $0xf4; BYTE $0xd5             // VPMULUDQ	Y5, Y3, Y10
	BYTE $0xc5; BYTE $0x5d; BYTE $0xf4; BYTE $0xdd             // VPMULUDQ	Y5, Y4, Y11

	VMOVDQU 160(DI), Y5

	BYTE    $0xc4; BYTE $0xe2; BYTE $0x15; BYTE $0x36; BYTE $0xed // VPERMD		Y5, Y13, Y5
	BYTE    $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y4, Y6
	VPADDQ  Y6, Y7, Y7
	VMOVDQU 32(DI), Y5
	BYTE    $0xc4; BYTE $0xe2; BYTE $0x15; BYTE $0x36; BYTE $0xed // VPERMD		Y5, Y13, Y5
	BYTE    $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y0, Y6
	VPADDQ  Y6, Y8, Y8
	BYTE    $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y1, Y6
	VPADDQ  Y6, Y9, Y9
	BYTE    $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y2, Y6
	VPADDQ  Y6, Y10, Y10
	BYTE    $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y3, Y6
	VPADDQ  Y6, Y11, Y11

	VMOVDQU 192(DI), Y5
	BYTE    $0xc4; BYTE $0xe2; BYTE $0x15; BYTE $0x36; BYTE $0xed // VPERMD		Y5, Y13, Y5
	BYTE    $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y3, Y6
	VPADDQ  Y6, Y7, Y7
	BYTE    $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y4, Y6
	VPADDQ  Y6, Y8, Y8

	VMOVDQU 64(DI), Y5
	BYTE    $0xc4; BYTE $0xe2; BYTE $0x15; BYTE $0x36; BYTE $0xed // VPERMD		Y5, Y13, Y5
	BYTE    $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y0, Y6
	VPADDQ  Y6, Y9, Y9
	BYTE    $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y1, Y6
	VPADDQ  Y6, Y10, Y10
	BYTE    $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y2, Y6
	VPADDQ  Y6, Y11, Y11

	VMOVDQU 224(DI), Y5
	BYTE    $0xc4; BYTE $0xe2; BYTE $0x15; BYTE $0x36; BYTE $0xed // VPERMD		Y5, Y13, Y5
	BYTE    $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y2, Y6
	VPADDQ  Y6, Y7, Y7
	BYTE    $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y3, Y6
	VPADDQ  Y6, Y8, Y8
	BYTE    $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y4, Y6
	VPADDQ  Y6, Y9, Y9

	VMOVDQU 96(DI), Y5
	BYTE    $0xc4; BYTE $0xe2; BYTE $0x15; BYTE $0x36; BYTE $0xed // VPERMD		Y5, Y13, Y5
	BYTE    $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y0, Y6
	VPADDQ  Y6, Y10, Y10
	BYTE    $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y1, Y6
	VPADDQ  Y6, Y11, Y11

	VMOVDQU 256(DI), Y5
	BYTE    $0xc4; BYTE $0xe2; BYTE $0x15; BYTE $0x36; BYTE $0xed // VPERMD		Y5, Y13, Y5
	BYTE    $0xc5; BYTE $0xf5; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y1,Y6
	VPADDQ  Y6, Y7, Y7
	BYTE    $0xc5; BYTE $0xed; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y2, Y6
	VPADDQ  Y6, Y8, Y8
	BYTE    $0xc5; BYTE $0xe5; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y3, Y6
	VPADDQ  Y6, Y9, Y9
	BYTE    $0xc5; BYTE $0xdd; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y4, Y6
	VPADDQ  Y6, Y10, Y10

	VMOVDQU 128(DI), Y5
	BYTE    $0xc4; BYTE $0xe2; BYTE $0x15; BYTE $0x36; BYTE $0xed // VPERMD		Y5, Y13, Y5
	BYTE    $0xc5; BYTE $0xfd; BYTE $0xf4; BYTE $0xf5             // VPMULUDQ	Y5, Y0, Y6
	VPADDQ  Y6, Y11, Y11

	VPSRLQ $26, Y10, Y5
	VPADDQ Y5, Y11, Y11
	VPAND  Y12, Y10, Y10
	VPSRLQ $26, Y11, Y5
	VPSLLQ $2, Y5, Y6
	VPADDQ Y6, Y5, Y5
	VPADDQ Y5, Y7, Y7
	VPAND  Y12, Y11, Y11
	VPSRLQ $26, Y7, Y5
	VPAND  Y12, Y7, Y0
	VPADDQ Y5, Y8, Y8
	VPSRLQ $26, Y8, Y5
	VPAND  Y12, Y8, Y1
	VPADDQ Y5, Y9, Y9
	VPSRLQ $26, Y9, Y5
	VPAND  Y12, Y9, Y2
	VPADDQ Y5, Y10, Y10
	VPSRLQ $26, Y10, Y5
	VPAND  Y12, Y10, Y3
	VPADDQ Y5, Y11, Y4

	VPSRLDQ $8, Y0, Y7
	VPSRLDQ $8, Y1, Y8
	VPSRLDQ $8, Y2, Y9
	VPSRLDQ $8, Y3, Y10
	VPSRLDQ $8, Y4, Y11

	VPADDQ Y7, Y0, Y0
	VPADDQ Y8, Y1, Y1
	VPADDQ Y9, Y2, Y2
	VPADDQ Y10, Y3, Y3
	VPADDQ Y11, Y4, Y4

	BYTE $0xc4; BYTE $0xe3; BYTE $0xfd; BYTE $0x00; BYTE $0xf8; BYTE $0xaa // VPERMQ	$170, Y0, Y7
	BYTE $0xc4; BYTE $0x63; BYTE $0xfd; BYTE $0x00; BYTE $0xc1; BYTE $0xaa // VPERMQ	$170, Y1, Y8
	BYTE $0xc4; BYTE $0x63; BYTE $0xfd; BYTE $0x00; BYTE $0xca; BYTE $0xaa // VPERMQ	$170, Y2, Y9
	BYTE $0xc4; BYTE $0x63; BYTE $0xfd; BYTE $0x00; BYTE $0xd3; BYTE $0xaa // VPERMQ	$170, Y3, Y10
	BYTE $0xc4; BYTE $0x63; BYTE $0xfd; BYTE $0x00; BYTE $0xdc; BYTE $0xaa // VPERMQ	$170, Y4, Y11

	VPADDQ Y7, Y0, Y0
	VPADDQ Y8, Y1, Y1
	VPADDQ Y9, Y2, Y2
	VPADDQ Y10, Y3, Y3
	VPADDQ Y11, Y4, Y4

DONE:
	VZEROUPPER
	MOVD X0, 304(DI)
	MOVD X1, 308(DI)
	MOVD X2, 312(DI)
	MOVD X3, 316(DI)
	MOVD X4, 320(DI)
	RET

TEXT ·finalizeAVX2(SB), $0-16
	MOVQ out+0(FP), SI
	MOVQ state+8(FP), DI

	VZEROUPPER

	BYTE $0xc5; BYTE $0xf9; BYTE $0x6e; BYTE $0x87; BYTE $0x30; BYTE $0x01; BYTE $0x00; BYTE $0x00 // VMOVD	304(DI), X0
	BYTE $0xc5; BYTE $0xf9; BYTE $0x6e; BYTE $0x8f; BYTE $0x34; BYTE $0x01; BYTE $0x00; BYTE $0x00 // VMOVD	308(DI), X1
	BYTE $0xc5; BYTE $0xf9; BYTE $0x6e; BYTE $0x97; BYTE $0x38; BYTE $0x01; BYTE $0x00; BYTE $0x00 // VMOVD	312(DI), X2
	BYTE $0xc5; BYTE $0xf9; BYTE $0x6e; BYTE $0x9f; BYTE $0x3c; BYTE $0x01; BYTE $0x00; BYTE $0x00 // VMOVD	316(DI), X3
	BYTE $0xc5; BYTE $0xf9; BYTE $0x6e; BYTE $0xa7; BYTE $0x40; BYTE $0x01; BYTE $0x00; BYTE $0x00 // VMOVD	320(DI), X4

	VMOVDQU addMaskAVX2<>(SB), X7

	VPSRLQ $26, X4, X5
	VPSLLQ $2, X5, X6
	VPADDQ X6, X5, X5
	VPADDQ X5, X0, X0
	VPAND  X7, X4, X4

	VPSRLQ $26, X0, X5
	VPAND  X7, X0, X0
	VPADDQ X5, X1, X1
	VPSRLQ $26, X1, X5
	VPAND  X7, X1, X1
	VPADDQ X5, X2, X2
	VPSRLQ $26, X2, X5
	VPAND  X7, X2, X2
	VPADDQ X5, X3, X3
	VPSRLQ $26, X3, X5
	VPAND  X7, X3, X3
	VPADDQ X5, X4, X4

	VPSLLQ $26, X1, X5
	VPXOR  X5, X0, X0
	VPSLLQ $52, X2, X5
	VPXOR  X5, X0, X0
	VPSRLQ $12, X2, X1
	VPSLLQ $14, X3, X5
	VPXOR  X5, X1, X1
	VPSLLQ $40, X4, X5
	VPXOR  X5, X1, X1

	VZEROUPPER

	MOVQ X0, AX
	MOVQ X1, BX

	ADDQ 288(DI), AX
	ADCQ 288+8(DI), BX
	MOVQ AX, (SI)
	MOVQ BX, 8(SI)

	RET
