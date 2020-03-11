import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TablePagination from "@material-ui/core/TablePagination";
import TableRow from "@material-ui/core/TableRow";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";
import IconButton from "@material-ui/core/IconButton";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBook, faCheckSquare } from "@fortawesome/free-solid-svg-icons";
import { faSlack, faGithub } from "@fortawesome/free-brands-svg-icons";

const getFullName = nameParts => {
  return nameParts;
};

const getAttribute = (row, key) => {
  if (!("Attributes" in row && key in row.Attributes)) {
    return null;
  }
  return row.Attributes[key];
};

function desc(a, b, orderBy) {
  if (getAttribute(b, orderBy) < getAttribute(a, orderBy)) {
    return -1;
  }
  if (getAttribute(b, orderBy) > getAttribute(a, orderBy)) {
    return 1;
  }
  return 0;
}

function stableSort(array, cmp) {
  const stabilizedThis = array.map((el, index) => [el, index]);
  stabilizedThis.sort((a, b) => {
    const order = cmp(a[0], b[0]);
    if (order !== 0) return order;
    return a[1] - b[1];
  });
  return stabilizedThis.map(el => el[0]);
}

function getSorting(order, orderBy) {
  return order === "desc"
    ? (a, b) => desc(a, b, orderBy)
    : (a, b) => -desc(a, b, orderBy);
}

const headCells = [
  { id: "name", numeric: false, disablePadding: false, label: "Service" },
  { id: "prod", numeric: true, disablePadding: false, label: "Prod" },
  { id: "owner", numeric: true, disablePadding: false, label: "Owner" },
  {
    id: "description",
    numeric: true,
    disablePadding: false,
    label: "Description"
  },
  { id: "type", numeric: true, disablePadding: false, label: "Type" },
  { id: "version", numeric: true, disablePadding: false, label: "Version" },
  { id: "link", numeric: true, disablePadding: false, label: "Link" }
];

function EnhancedTableHead(props) {
  const { classes, order, orderBy, onRequestSort } = props;
  const createSortHandler = property => event => {
    onRequestSort(event, property);
  };

  return (
    <TableHead>
      <TableRow>
        {headCells.map(headCell => (
          <TableCell
            key={headCell.id}
            align={headCell.numeric ? "right" : "left"}
            padding={headCell.disablePadding ? "none" : "default"}
            sortDirection={orderBy === headCell.id ? order : false}
          >
            <TableSortLabel
              active={orderBy === headCell.id}
              direction={orderBy === headCell.id ? order : "asc"}
              onClick={createSortHandler(headCell.id)}
            >
              {headCell.label}
              {orderBy === headCell.id ? (
                <span className={classes.visuallyHidden}>
                  {order === "desc" ? "sorted descending" : "sorted ascending"}
                </span>
              ) : null}
            </TableSortLabel>
          </TableCell>
        ))}
      </TableRow>
    </TableHead>
  );
}

const useStyles = makeStyles(theme => ({
  root: {
    width: "100%"
  },
  paper: {
    width: "100%",
    marginBottom: theme.spacing(2),
    padding: "50px"
  },
  table: {
    minWidth: 750
  },
  visuallyHidden: {
    border: 0,
    clip: "rect(0 0 0 0)",
    height: 1,
    margin: -1,
    overflow: "hidden",
    padding: 0,
    position: "absolute",
    top: 20,
    width: 1
  }
}));

function ServiceListView(props) {
  const classes = useStyles();
  const rows = props.items;
  const [order, setOrder] = React.useState("asc");
  const [orderBy, setOrderBy] = React.useState("service");
  const [selected, setSelected] = React.useState([]);
  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(5);

  const handleRequestSort = (event, property) => {
    const isAsc = orderBy === property && order === "asc";
    setOrder(isAsc ? "desc" : "asc");
    setOrderBy(property);
  };

  const handleSelectAllClick = event => {
    if (event.target.checked) {
      const newSelecteds = rows.map(n => getFullName(n.Name));
      setSelected(newSelecteds);
      return;
    }
    setSelected([]);
  };

  const handleClick = (event, name) => {
    const selectedIndex = selected.indexOf(name);
    let newSelected = [];

    if (selectedIndex === -1) {
      newSelected = newSelected.concat(selected, name);
    } else if (selectedIndex === 0) {
      newSelected = newSelected.concat(selected.slice(1));
    } else if (selectedIndex === selected.length - 1) {
      newSelected = newSelected.concat(selected.slice(0, -1));
    } else if (selectedIndex > 0) {
      newSelected = newSelected.concat(
        selected.slice(0, selectedIndex),
        selected.slice(selectedIndex + 1)
      );
    }

    setSelected(newSelected);
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = event => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const isSelected = name => selected.indexOf(name) !== -1;

  const emptyRows =
    rowsPerPage - Math.min(rowsPerPage, rows.length - page * rowsPerPage);

  return (
    <div className={classes.root}>
      <Paper className={classes.paper}>
        <Typography className={classes.title} variant="h6" id="tableTitle">
          Services
        </Typography>
        <TableContainer>
          <Table
            className={classes.table}
            aria-labelledby="tableTitle"
            size="small"
            aria-label="enhanced table"
          >
            <EnhancedTableHead
              classes={classes}
              numSelected={selected.length}
              order={order}
              orderBy={orderBy}
              onSelectAllClick={handleSelectAllClick}
              onRequestSort={handleRequestSort}
              rowCount={rows.length}
            />
            <TableBody>
              {stableSort(rows, getSorting(order, orderBy))
                .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                .map((row, index) => {
                  const isItemSelected = isSelected(getFullName(row.Name));
                  const labelId = `enhanced-table-checkbox-${index}`;
                  return (
                    <TableRow
                      hover
                      onClick={event =>
                        handleClick(event, getFullName(row.Name))
                      }
                      role="checkbox"
                      aria-checked={isItemSelected}
                      tabIndex={-1}
                      key={getFullName(row.Name)}
                      selected={isItemSelected}
                    >
                      <TableCell
                        component="th"
                        id={labelId}
                        scope="row"
                        padding="none"
                      >
                        {getFullName(row.Name)}
                      </TableCell>
                      <TableCell align="right">
                        <a href={window.location.origin + row.Path}>
                          <IconButton>
                            <FontAwesomeIcon icon={faCheckSquare} />
                          </IconButton>
                        </a>
                      </TableCell>
                      <TableCell align="right">
                        {getAttribute(row, "owner.name")}
                      </TableCell>
                      <TableCell align="right">
                        {getAttribute(row, "description")}
                      </TableCell>
                      <TableCell align="right">
                        {row.Type}
                      </TableCell>
                      <TableCell align="right">
                        {getAttribute(row, "version")}
                      </TableCell>
                      <TableCell align="right">
                        <a href={getAttribute(row, "docs.url")}>
                          <IconButton>
                            <FontAwesomeIcon icon={faBook} />
                          </IconButton>
                        </a>
                        <a href={getAttribute(row, "team.slack")}>
                          <IconButton>
                            <FontAwesomeIcon icon={faSlack} />
                          </IconButton>
                        </a>
                        <a href={getAttribute(row, "repo.url")}>
                          <IconButton>
                            <FontAwesomeIcon icon={faGithub} />
                          </IconButton>
                        </a>
                      </TableCell>
                    </TableRow>
                  );
                })}
              {emptyRows > 0 && (
                <TableRow style={{ height: 33 * emptyRows }}>
                  <TableCell colSpan={6} />
                </TableRow>
              )}
            </TableBody>
          </Table>
        </TableContainer>
        <TablePagination
          rowsPerPageOptions={[5, 10, 25]}
          component="div"
          count={rows.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onChangePage={handleChangePage}
          onChangeRowsPerPage={handleChangeRowsPerPage}
        />
      </Paper>
    </div>
  );
}

export default ServiceListView;
