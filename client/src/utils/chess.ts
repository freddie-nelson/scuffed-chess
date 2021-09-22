export enum Color {
  Black,
  White,
}

export enum Class {
  Queen,
  King,
  Rook,
  Bishop,
  Knight,
  Pawn,
}

export interface Piece {
  color: Color;
  class: Class;
}

export interface Spot {
  piece?: Piece;
  containsPiece: boolean;
  file: number;
  rank: number;
}
export interface Game {
  color: Color;
  opponentColor: Color;
  board: Spot[][];
  ended: boolean;
  endState: string;
  turn: Color;
  halfmoves: number;
  fullmoves: number;
  code: string;
}

export default class Chess {
  static fromFENString(fen: string): Game {
    const g: any = {};

    const piecePlacements = fen.split("/");

    const last = piecePlacements[7].split(" ");
    piecePlacements[7] = last[0];
    const fields = last.slice(1);

    // current turn
    if (fields[0] === "b") {
      g.turn = Color.Black;
    } else {
      g.turn = Color.White;
    }

    // castling rights
    g.blackCastling = { kingside: false, queenside: false };
    g.whiteCastling = { kingside: false, queenside: false };

    for (const rights of fields[1]) {
      if (rights.toLowerCase() === rights) {
        if (rights === "k") {
          g.blackCastling.kingside = true;
        } else if (rights === "q") {
          g.blackCastling.queenside = true;
        }
      } else {
        if (rights === "K") {
          g.whiteCastling.kingside = true;
        } else if (rights === "Q") {
          g.whiteCastling.queenside = true;
        }
      }
    }

    // en passant targets
    // if (fields[2] !== "-") {
    // 	const [file, rank] = this.locationToFileAndRank(fields[2])
    // 	g.board.grid[file][rank].passantTarget = 2
    // }

    // fullmoves and halfmoves
    g.halfmoves = Number(fields[3]);
    g.fullmoves = Number(fields[4]);

    // place pieces
    g.board = [];
    for (let file = 0; file < 8; file++) {
      g.board.push([]);

      for (let rank = 0; rank < 8; rank++) {
        const spot: Spot = {
          containsPiece: false,
          file: file,
          rank: rank,
        };

        g.board[file].push(spot);
      }
    }

    for (let rank = 0; rank < piecePlacements.length; rank++) {
      const fenRank = piecePlacements[rank];
      let file = 0;

      for (const char of fenRank) {
        let color: Color;
        let pieceClass: Class = Class.Pawn;

        if (!isNaN(Number(char))) {
          file += Number(char);
          continue;
        } else if (char.toLowerCase() === char) {
          color = Color.Black;
        } else {
          color = Color.White;
        }

        switch (char.toUpperCase()) {
          case "Q":
            pieceClass = Class.Queen;
            break;
          case "K":
            pieceClass = Class.King;
            break;
          case "R":
            pieceClass = Class.Rook;
            break;
          case "B":
            pieceClass = Class.Bishop;
            break;
          case "N":
            pieceClass = Class.Knight;
            break;
          case "P":
            pieceClass = Class.Pawn;
            break;
        }

        g.board[file][rank].containsPiece = true;
        g.board[file][rank].piece = { color: color, class: pieceClass };
        file++;
      }
    }

    console.log(g);

    return g;
  }

  static locationToFileAndRank(loc: string): number[] {
    const file = loc.charCodeAt(0) - "a".charCodeAt(0);
    const rank = 8 - loc.charCodeAt(1) - "0".charCodeAt(0);

    return [file, rank];
  }
}

// window.Chess = Chess;
